package services

import (
	"context"

	"github.com/micro-eshop/catalog/internal/core/model"
)

type CatalogService interface {
	GetProductById(ctx context.Context, id model.ProductId) (*model.Product, error)
}
