package mongo

import (
	"context"

	"github.com/micro-eshop/catalog/internal/core/model"
)

type mongoCatalogRepository struct {
}

func NewMongoCatalogRepository() *mongoCatalogRepository {
	return &mongoCatalogRepository{}
}

func (r *mongoCatalogRepository) GetProductById(ctx context.Context, id model.ProductId) (*model.Product, error) {
	return nil, nil
}

func (r *mongoCatalogRepository) GetProductByIds(ctx context.Context, ids []model.ProductId) ([]*model.Product, error) {
	return nil, nil
}
