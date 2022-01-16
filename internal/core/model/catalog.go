package model

import "errors"

type ProductId int

type Product struct {
	Id ProductId
}

func NewProduct(id ProductId) *Product {
	return &Product{
		Id: id,
	}
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
