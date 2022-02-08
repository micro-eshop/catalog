package cmd

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"github.com/micro-eshop/catalog/internal/core/services"
	"github.com/micro-eshop/catalog/internal/core/usecase"
	"github.com/micro-eshop/catalog/internal/infrastructure/repositories"
	log "github.com/sirupsen/logrus"
)

type ImportProductsCmd struct {
	mongoConnection string
}

func (*ImportProductsCmd) Name() string     { return "run-import" }
func (*ImportProductsCmd) Synopsis() string { return "Import all products" }
func (*ImportProductsCmd) Usage() string {
	return ""
}

func (p *ImportProductsCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.mongoConnection, "mongoConnection", "mongodb://localhost:27017", "mongodb connection string")
}

func (p *ImportProductsCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	mongodbClient, err := repositories.NewClient(ctx, p.mongoConnection)
	if err != nil {
		log.WithError(err).Error("can't create mongodb client")
		return subcommands.ExitFailure
	}
	defer mongodbClient.Close(ctx)
	repo := repositories.NewMongoCatalogRepository(mongodbClient)
	service := services.NewCatalogImportService(repo)
	importUc := usecase.NewImportProductsUseCase(service, services.NewProductsSourceDataProvider(), nil)
	return subcommands.ExitSuccess
}
