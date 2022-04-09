package repositories

import (
	"context"

	"github.com/micro-eshop/catalog/pkg/core/model"
)

type ProductSearchParams struct {
	Name        string
	Brand       string
	PriceFrom   float64
	PriceTo     float64
	InPromotion bool
}

type CatalogReader interface {
	GetProductById(ctx context.Context, id model.ProductId) (*model.Product, error)
	GetProductByIds(ctx context.Context, ids ...model.ProductId) ([]*model.Product, error)
	Search(ctx context.Context, params ProductSearchParams) ([]*model.Product, error)
}

type CatalogWriter interface {
	Insert(ctx context.Context, product *model.Product) error
}

type CatalogRepository interface {
	CatalogReader
	CatalogWriter
}
