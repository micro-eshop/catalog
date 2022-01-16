package repositories

import (
	"context"

	"github.com/micro-eshop/catalog/internal/core/model"
)

type CatalogRepository interface {
	GetProductById(ctx context.Context, id model.ProductId) (*model.Product, error)
	GetProductByIds(ctx context.Context, ids []model.ProductId) ([]*model.Product, error)
}
