package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/micro-eshop/catalog/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var v *viper.Viper

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
	v = viper.New()
	v.SetDefault("POSTGRES_CONNECTION", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
}

func main() {

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(cmd.NewRunApiCmd(v), "")
	subcommands.Register(&cmd.ImportProductsCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
