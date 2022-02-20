package cmd

import (
	"context"
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/google/subcommands"
	"github.com/micro-eshop/catalog/internal/env"
	"github.com/micro-eshop/catalog/internal/nats"
	"github.com/micro-eshop/catalog/internal/postgres"
	"github.com/micro-eshop/catalog/pkg/core/services"
	"github.com/micro-eshop/catalog/pkg/core/usecase"
	"github.com/micro-eshop/catalog/pkg/handlers"
	log "github.com/sirupsen/logrus"
)

type RunApiCmd struct {
	addr         string
	postgresConn string
}

func (*RunApiCmd) Name() string     { return "run-api" }
func (*RunApiCmd) Synopsis() string { return "Run catalog api" }
func (*RunApiCmd) Usage() string {
	return `run-api`
}

func (p *RunApiCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.addr, "addr", ":8080", "address to listen")
	f.StringVar(&p.postgresConn, "postgresConn", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "postgresConn connection string")
}

func (p *RunApiCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	r := gin.Default()
	log.Infoln("Start import products")
	postgresClient, err := postgres.NewPostgresClient(ctx, p.postgresConn)
	if err != nil {
		log.WithError(err).Error("can't create postgres  client")
		return subcommands.ExitFailure
	}
	defer postgresClient.Close(ctx)
	publisher, err := nats.NewPublisher(env.GetEnvOrDefault("NATS_URL", "nats://nats:4222"))
	if err != nil {
		log.WithError(err).Error("can't create publisher")
		return subcommands.ExitFailure
	}
	defer publisher.Close()

	repo := postgres.NewPostgresCatalogRepository(postgresClient)
	service := services.NewCatalogService(repo)

	getById := usecase.NewGetProductByIdUseCase(service)
	getByIds := usecase.NewGetProductByIdsUseCase(service)
	catalog := handlers.NewCatalogHandler(getById, getByIds)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	catalog.Setup(r)
	if err := r.Run(p.addr); err != nil {
		log.WithError(err).WithContext(ctx).Errorln("failed to run api")
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}
