package repositories

import (
	"context"

	"github.com/micro-eshop/catalog/internal/core/model"
)

type CatalogReader interface {
	GetProductById(ctx context.Context, id model.ProductId) (*model.Product, error)
	GetProductByIds(ctx context.Context, ids []model.ProductId) ([]*model.Product, error)
}

type CatalogWriter interface {
	Insert(ctx context.Context, product *model.Product) error
}

type CatalogRepository interface {
	CatalogReader
	CatalogWriter
}
