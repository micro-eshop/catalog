package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/micro-eshop/catalog/cmd"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&cmd.RunApiCmd{}, "")
	subcommands.Register(&cmd.ImportProductsCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
