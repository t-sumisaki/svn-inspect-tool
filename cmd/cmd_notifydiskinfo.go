package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"t-sumisaki/svn-inspect-tool/diskutil"

	"github.com/google/subcommands"
	"github.com/slack-go/slack"
)

type NotifyDiskInfoConfig struct {
	TargetMountPoint string `yaml:"TargetMountPoint"`
	ScanTargetPath   string `yaml:"ScanTargetPath"`
	SlackWebhookURL  string `yaml:"SlackWebhookURL"`
}

type notifyDiskInfoCmd struct {
	dryrun bool
}

func (*notifyDiskInfoCmd) Name() string     { return "notifydiskinfo" }
func (*notifyDiskInfoCmd) Synopsis() string { return "Post Disk information to slack" }
func (*notifyDiskInfoCmd) Usage() string {
	return `notifydiskinfo [-dryrun]
	Post Disk information to slack`
}

func (c *notifyDiskInfoCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.dryrun, "dryrun", false, "run without slack api")
}

func (c *notifyDiskInfoCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

	logfile.Info().Bool("dryrun", c.dryrun).Msg("NotifyDiskInfo command start")

	var _config ConfigFile
	if err := loadConfig(&_config); err != nil {
		logfile.Error().Err(err).Msg("loadConfig failed.")
		return subcommands.ExitFailure
	}

	conf := _config.NotifyDiskInfo

	dfResult, err := diskutil.GetDFResult()
	if err != nil {
		logfile.Error().Err(err).Msg("failed to df command")
		return subcommands.ExitFailure
	}

	targetMount := diskutil.FindByMount(dfResult, conf.TargetMountPoint)
	if targetMount == nil {
		logfile.Error().Str("mountpoint", conf.TargetMountPoint).Msg("target mountpoint is not found.")
		return subcommands.ExitFailure
	}

	duResult, err := diskutil.GetDUResult(conf.ScanTargetPath)
	if err != nil {
		logfile.Error().Err(err).Msg("failed to du command")
		return subcommands.ExitFailure
	}

	var builder strings.Builder
	for _, entry := range duResult {
		builder.WriteString(fmt.Sprintf("%s\t\t%s", entry.Size, entry.Path))
	}

	if c.dryrun {
		fmt.Println("Result:")
		fmt.Println(dfResult)
		fmt.Println(duResult)
		return subcommands.ExitSuccess
	}

	bm := &slack.WebhookMessage{
		Text: "*SVNディスク使用量レポート*",
		Blocks: &slack.Blocks{
			BlockSet: []slack.Block{
				slack.NewSectionBlock(&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: "*SVNディスク使用量レポート*",
				}, nil, nil),
				slack.NewSectionBlock(&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: fmt.Sprintf("*[%s]空き容量:* %s (%s) (%s / %s)",
						targetMount.MountedOn,
						targetMount.Available,
						targetMount.UsePercent,
						targetMount.Used,
						targetMount.Size),
				}, nil, nil),
				slack.NewSectionBlock(&slack.TextBlockObject{
					Type: slack.MarkdownType,
					Text: fmt.Sprintf("```%s```", builder.String()),
				}, nil, nil),
			},
		},
	}

	logfile.Info().Str("url", conf.SlackWebhookURL).Msg("Send slack webhook")
	if err := slack.PostWebhook(conf.SlackWebhookURL, bm); err != nil {
		logfile.Err(err).Msg("Slack postwebhook failed")
	}

	logfile.Info().Msg("NotifyLockInfo command completed")
	return subcommands.ExitSuccess

}
