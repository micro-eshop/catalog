package services

import (
	"context"

	"github.com/micro-eshop/catalog/internal/core/model"
	"github.com/micro-eshop/catalog/internal/core/repositories"
)

type CatalogService interface {
	GetProductById(ctx context.Context, id model.ProductId) (*model.Product, error)
	GetProductByIds(ctx context.Context, ids []model.ProductId) ([]*model.Product, error)
}

type catalogService struct {
	repo repositories.CatalogRepository
}

func NewCatalogService(repo repositories.CatalogRepository) CatalogService {
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
	repo repositories.CatalogRepository
}

func NewCatalogImportService(repo repositories.CatalogRepository) CatalogService {
	return &catalogService{
		repo: repo,
	}
}

func (s *catalogImportService) Store(ctx context.Context, products []*model.Product) error {
	return nil
}
