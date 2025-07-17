package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/rs/zerolog"
)

var (
	logfile zerolog.Logger
)

func main() {

	file, err := os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logfile = zerolog.New(file).With().Timestamp().Logger()

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&nofityLockInfoCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
