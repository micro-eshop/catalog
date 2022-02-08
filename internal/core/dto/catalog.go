package dto

import (
	"github.com/micro-eshop/catalog/internal/core/model"
)

type ProductDto struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Brand          string   `json:"brand"`
	Description    string   `json:"description"`
	Price          float64  `json:"price"`
	PromotionPrice *float64 `json:"promotionPrice"`
}

func NewProductDto(product *model.Product) *ProductDto {
	return &ProductDto{ID: int(product.ID), Name: product.Name, Brand: product.Brand, Description: product.Description, Price: product.Price, PromotionPrice: product.PromotionPrice}
}
