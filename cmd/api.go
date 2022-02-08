package cmd

import (
	"context"
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/google/subcommands"
	log "github.com/sirupsen/logrus"
)

type RunApiCmd struct {
	addr string
}

func (*RunApiCmd) Name() string     { return "run-api" }
func (*RunApiCmd) Synopsis() string { return "Run catalog api" }
func (*RunApiCmd) Usage() string {
	return `run-api`
}

func (p *RunApiCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.addr, "addr", ":8080", "address to listen")
}

func (p *RunApiCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(p.addr); err != nil {
		log.WithError(err).WithContext(ctx).Errorln("failed to run api")
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
