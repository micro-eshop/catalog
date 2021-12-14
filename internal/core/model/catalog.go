package model

type ProductId int

type Product struct {
	Id ProductId
}

func NewProduct(id ProductId) *Product {
	return &Product{
		Id: id,
	}
}
