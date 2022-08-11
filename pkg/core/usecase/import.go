package usecase

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/micro-eshop/catalog/pkg/core/services"
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
	data := uc.source.Provide(ctx)
	stream := uc.service.Store(ctx, data)
	var err error
	for product := range stream {
		publishErr := uc.publisher.PublishProductCreated(ctx, services.NewProductCreated(product))
		if publishErr != nil {
			err = multierror.Append(err, publishErr)
		}
	}
	return err
}
