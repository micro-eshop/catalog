package usecase

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/micro-eshop/catalog/internal/core/model"
	"github.com/micro-eshop/catalog/internal/core/services"
)

type getProductByIdUseCase struct {
	service services.CatalogService
}

func NewGetProductByIdUseCase(service services.CatalogService) *getProductByIdUseCase {
	return &getProductByIdUseCase{
		service: service,
	}
}

func (uc *getProductByIdUseCase) Execute(ctx context.Context, id model.ProductId) (*model.Product, error) {
	return uc.service.GetProductById(ctx, id)
}

type getProductByIdsUseCase struct {
	service services.CatalogService
}

func NewGetProductByIdsUseCase(service services.CatalogService) *getProductByIdsUseCase {
	return &getProductByIdsUseCase{
		service: service,
	}
}

func (uc *getProductByIdsUseCase) Execute(ctx context.Context, ids []model.ProductId) ([]*model.Product, error) {
	return uc.service.GetProductByIds(ctx, ids)
}

type importProductsUseCase struct {
	service   services.CatalogImportService
	source    services.ProductsSourceDataProvider
	publisher services.ProductCreatedPublisher
}

func NewImportProductsUseCase(service services.CatalogImportService, source services.ProductsSourceDataProvider) *importProductsUseCase {
	return &importProductsUseCase{
		service: service,
		source:  source,
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
