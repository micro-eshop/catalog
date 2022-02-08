package services

import (
	"context"

	"github.com/micro-eshop/catalog/internal/core/model"
	"github.com/micro-eshop/catalog/internal/core/repositories"
)

type CatalogService interface {
	GetProductById(ctx context.Context, id model.ProductId) (*model.Product, error)
	GetProductByIds(ctx context.Context, ids []model.ProductId) ([]*model.Product, error)
	Search(ctx context.Context, params repositories.ProductSearchParams) ([]*model.Product, error)
}

type catalogService struct {
	repo repositories.CatalogReader
}

func NewCatalogService(repo repositories.CatalogReader) *catalogService {
	return &catalogService{
		repo: repo,
	}
}

func (s *catalogService) GetProductById(ctx context.Context, id model.ProductId) (*model.Product, error) {
	err := model.ValidateProductId(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetProductById(ctx, id)
}

func (s *catalogService) GetProductByIds(ctx context.Context, ids []model.ProductId) ([]*model.Product, error) {
	err := model.ValidateProductIds(ids)
	if err != nil {
		return nil, err
	}
	return s.repo.GetProductByIds(ctx, ids)
}

type CatalogImportService interface {
	Store(ctx context.Context, products []*model.Product) error
}

type catalogImportService struct {
	repo repositories.CatalogWriter
}

func NewCatalogImportService(repo repositories.CatalogRepository) *catalogImportService {
	return &catalogImportService{
		repo: repo,
	}
}

func (s *catalogImportService) Store(ctx context.Context, products []*model.Product) error {
	for _, product := range products {
		err := model.ValidateProduct(product)
		if err != nil {
			return err
		}
		err = s.repo.Insert(ctx, product)
		if err != nil {
			return err
		}
	}
	return nil
}

type ProductsSourceDataProvider interface {
	Provide(ctx context.Context) ([]*model.Product, error)
}

type productsSourceDataProvider struct {
}

func NewProductsSourceDataProvider() *productsSourceDataProvider {
	return &productsSourceDataProvider{}
}

func (s *productsSourceDataProvider) Provide(ctx context.Context) ([]*model.Product, error) {
	return make([]*model.Product, 0), nil
}
