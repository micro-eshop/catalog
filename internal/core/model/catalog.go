package model

import "errors"

type ProductId int

type Product struct {
	ID             ProductId
	Name           string
	Brand          string
	Description    string
	Price          float64
	PromotionPrice *float64
}

type Products = []Product

func NewProduct(id ProductId, name, brand, description string, price float64, promotionPrice *float64) *Product {
	return &Product{ID: id, Name: name, Brand: brand, Description: description, Price: price, PromotionPrice: promotionPrice}
}

func ValidateProductId(id ProductId) error {
	if id > 0 {
		return nil
	}
	return errors.New("ProductId must be greater than 0")
}

func ValidateProductIds(ids []ProductId) error {
	for id := range ids {
		err := ValidateProductId(ProductId(id))
		if err != nil {
			return err
		}
	}
	return nil
}
