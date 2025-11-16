package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"t-sumisaki/svn-inspect-tool/svnadmin"

	"github.com/google/subcommands"

	"github.com/slack-go/slack"
)

type NotifyLockInfoConfig struct {
	Name            string `yaml:"Name"`
	RepositoryPath  string `yaml:"RepositoryPath"`
	SlackWebhookURL string `yaml:"SlackWebhookURL"`
}

type NotifyLockInfoConfigSet map[string]NotifyLockInfoConfig

type nofityLockInfoCmd struct {
	profile string
	dryrun  bool
}

func (*nofityLockInfoCmd) Name() string     { return "notifylockinfo" }
func (*nofityLockInfoCmd) Synopsis() string { return "Post SVN lock status message to Slack" }
func (*nofityLockInfoCmd) Usage() string {
	return `notifylockinfo -profile
	Post SVN lock status message to Slack`
}

func (c *nofityLockInfoCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.profile, "profile", "", "config profile")
	f.BoolVar(&c.dryrun, "dryrun", false, "run without slack api")
}

func (c *nofityLockInfoCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	logfile.Info().Str("profile", c.profile).Bool("dryrun", c.dryrun).Msg("NotifyLockInfo command start")
	var config ConfigFile

	if err := loadConfig(&config); err != nil {
		logfile.Err(err).Msg("loadConfig Failed.")
		return subcommands.ExitFailure
	}

	profile, ok := config.NotifyLockInfo[c.profile]
	if !ok {
		logfile.Error().Str("profile", c.profile).Msg("profile not found")
		return subcommands.ExitFailure
	}

	if profile.RepositoryPath == "" {
		logfile.Error().Str("profile", c.profile).Msg("RepositoryPath is not defined")
		return subcommands.ExitFailure
	}

	if profile.SlackWebhookURL == "" {
		logfile.Error().Str("profile", c.profile).Msg("SlackWebhookURL is not defined")
		return subcommands.ExitFailure
	}

	logfile.Info().Str("Name", profile.Name).Str("RepositoryPath", profile.RepositoryPath).Msg("Start Query")

	lockinfo, err := svnadmin.GetLockInfo(profile.RepositoryPath)
	if err != nil {
		logfile.Err(err).Msg("failed to svnadmin command")
		return subcommands.ExitFailure
	}

	grouped := make(map[string][]svnadmin.SvnLockInfo)
	for _, lock := range lockinfo {
		grouped[lock.Owner] = append(grouped[lock.Owner], lock)
	}

	var builder strings.Builder
	for owner, infos := range grouped {
		builder.WriteString(fmt.Sprintf("User: %s\n", owner))
		for _, info := range infos {
			builder.WriteString(fmt.Sprintf("%s\n", info.Path))
		}
		builder.WriteString("\n")
	}

	if builder.Len() == 0 {
		builder.WriteString("No locked asset")
	}

	if c.dryrun {
		fmt.Printf("Result: \n%s\n", builder.String())
		return subcommands.ExitSuccess
	}

	bm := &slack.WebhookMessage{
		Text: fmt.Sprintf("*SVNロック情報* (%s)", profile.Name),
		Blocks: &slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewSectionBlock(&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: fmt.Sprintf("*SVNロック情報* (%s)", profile.Name),
				}, nil, nil),
				slack.NewSectionBlock(&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: fmt.Sprintf("```%s```", builder.String()),
				}, nil, nil),
			},
		},
	}

	logfile.Info().Str("url", profile.SlackWebhookURL).Msg("Send slack webhook")
	if err := slack.PostWebhook(profile.SlackWebhookURL, bm); err != nil {
		logfile.Err(err).Msg("Slack postwebhook failed")
	}

	logfile.Info().Str("profile", profile.Name).Msg("NotifyLockInfo command completed")
	return subcommands.ExitSuccess
}
