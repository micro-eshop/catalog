package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/micro-eshop/catalog/cmd"
)

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
