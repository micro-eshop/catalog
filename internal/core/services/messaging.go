package services

import (
	"context"

	"github.com/micro-eshop/catalog/internal/core/model"
)

type ProductCreated struct {
}

func NewProductCreated(p *model.Product) *ProductCreated {
	return &ProductCreated{}
}

type ProductCreatedPublisher interface {
	Publish(ctx context.Context, event *ProductCreated) error
}
