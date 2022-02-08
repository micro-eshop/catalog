package usecase

import (
	"context"

	"github.com/micro-eshop/catalog/internal/core/dto"
	"github.com/micro-eshop/catalog/internal/core/model"
	"github.com/micro-eshop/catalog/internal/core/repositories"
	"github.com/micro-eshop/catalog/internal/core/services"
)

type GetProductByIdUseCase struct {
	service services.CatalogService
}

func NewGetProductByIdUseCase(service services.CatalogService) *GetProductByIdUseCase {
	return &GetProductByIdUseCase{
		service: service,
	}
}

func (uc *GetProductByIdUseCase) Execute(ctx context.Context, id model.ProductId) (*dto.ProductDto, error) {
	product, err := uc.service.GetProductById(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, nil
	}
	return dto.NewProductDto(product), nil
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

type searchProductsUseCase struct {
	service services.CatalogService
}

func NewSearchProductsUseCase(service services.CatalogService) *getProductByIdsUseCase {
	return &getProductByIdsUseCase{
		service: service,
	}
}

func (uc *searchProductsUseCase) Execute(ctx context.Context, params repositories.ProductSearchParams) ([]*model.Product, error) {
	return uc.service.Search(ctx, params)
}
