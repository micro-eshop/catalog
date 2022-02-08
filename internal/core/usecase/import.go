package usecase

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/micro-eshop/catalog/internal/core/services"
)

type importProductsUseCase struct {
	service   services.CatalogImportService
	source    services.ProductsSourceDataProvider
	publisher services.ProductCreatedPublisher
}

func NewImportProductsUseCase(service services.CatalogImportService, source services.ProductsSourceDataProvider, publisher services.ProductCreatedPublisher) *importProductsUseCase {
	return &importProductsUseCase{
		service:   service,
		source:    source,
		publisher: publisher,
	}
}

func (uc *importProductsUseCase) Execute(ctx context.Context) error {
	data, err := uc.source.Provide(ctx)
	if err != nil {
		return err
	}
	err = uc.service.Store(ctx, data)
	if err != nil {
		return err
	}
	for _, product := range data {
		publishErr := uc.publisher.Publish(ctx, services.NewProductCreated(product))
		if publishErr != nil {
			err = multierror.Append(err, publishErr)
		}
	}
	return err
}
