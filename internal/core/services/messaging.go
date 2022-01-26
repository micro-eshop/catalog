package services

import (
	"context"

	"github.com/micro-eshop/catalog/internal/core/model"
)

type ProductCreated struct {
	ID             int
	Name           string
	Brand          string
	Description    string
	Price          float64
	PromotionPrice *float64
}

func NewProductCreated(p *model.Product) *ProductCreated {
	return &ProductCreated{ID: int(p.ID), Name: p.Name, Brand: p.Brand, Description: p.Description, Price: p.Price, PromotionPrice: p.PromotionPrice}
}

type ProductCreatedPublisher interface {
	Publish(ctx context.Context, event *ProductCreated) error
}
