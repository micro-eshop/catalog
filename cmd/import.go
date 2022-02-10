package cmd

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/micro-eshop/catalog/internal/common/env"
	"github.com/micro-eshop/catalog/internal/core/services"
	"github.com/micro-eshop/catalog/internal/core/usecase"
	"github.com/micro-eshop/catalog/internal/infrastructure/messaging"
	"github.com/micro-eshop/catalog/internal/infrastructure/repositories"
	log "github.com/sirupsen/logrus"
)

type ImportProductsCmd struct {
	postgresConn string
}

func (*ImportProductsCmd) Name() string     { return "run-import" }
func (*ImportProductsCmd) Synopsis() string { return "Import all products" }
func (*ImportProductsCmd) Usage() string {
	return ""
}

func (p *ImportProductsCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.postgresConn, "postgresConn", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable", "postgresConn connection string")
}

func (p *ImportProductsCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	log.Infoln("Start import products")
	postgresClient, err := repositories.NewPostgresClient(ctx, p.postgresConn)
	if err != nil {
		log.WithError(err).Error("can't create postgres  client")
		return subcommands.ExitFailure
	}
	defer postgresClient.Close(ctx)
	publisher, err := messaging.NewPublisher(env.GetEnvOrDefault("NATS_URL", "nats://nats:4222"))
	if err != nil {
		log.WithError(err).Error("can't create publisher")
		return subcommands.ExitFailure
	}
	defer publisher.Close()

	repo := repositories.NewPostgresCatalogRepository(postgresClient)
	service := services.NewCatalogImportService(repo)

	importUc := usecase.NewImportProductsUseCase(service, services.NewProductsSourceDataProvider(), publisher)

	err = importUc.Execute(ctx)
	if err != nil {
		log.WithError(err).Error("can't import products")
		return subcommands.ExitFailure
	}
	log.Infoln("Finish import products")
	return subcommands.ExitSuccess
}
