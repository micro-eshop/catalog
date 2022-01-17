package usecase

import (
	"context"

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
