package services

import (
	"context"

	"github.com/micro-eshop/catalog/internal/core/model"
)

type ProductCreated struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Brand          string   `json:"brand"`
	Description    string   `json:"description"`
	Price          float64  `json:"price"`
	PromotionPrice *float64 `json:"promotion_price,omitempty"`
}

func NewProductCreated(p *model.Product) *ProductCreated {
	return &ProductCreated{ID: int(p.ID), Name: p.Name, Brand: p.Brand, Description: p.Description, Price: p.Price, PromotionPrice: p.PromotionPrice}
}

type ProductCreatedPublisher interface {
	Publish(ctx context.Context, event *ProductCreated) error
}
