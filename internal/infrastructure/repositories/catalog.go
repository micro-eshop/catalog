package repositories

import (
	"github.com/micro-eshop/catalog/internal/core/model"
	"go.mongodb.org/mongo-driver/bson"
)

const productsCollectionName = "products"
const catalogDatabase = "catalog"

type mongoProduct struct {
	ProductID      int      `json:"_id,omitempty" bson:"_id,omitempty"`
	Name           string   `json:"name,omitempty" bson:"name,omitempty"`
	Brand          string   `json:"brand,omitempty" bson:"brand,omitempty"`
	Description    string   `json:"description,omitempty" bson:"description,omitempty"`
	Price          float64  `json:"price,omitempty" bson:"price,omitempty"`
	PromotionPrice *float64 `json:"promotionPrice,omitempty" bson:"promotionPrice,omitempty"`
}

func newMongoProduct(product *model.Product) *mongoProduct {
	return &mongoProduct{ProductID: int(product.ID), Name: product.Name, Brand: product.Brand, Description: product.Description, Price: product.Price, PromotionPrice: product.PromotionPrice}
}

func toInsertMongoDocument(product *mongoProduct) bson.M {
	p := bson.M{
		"_id":            product.ProductID,
		"name":           product.Name,
		"brand":          product.Brand,
		"description":    product.Description,
		"price":          product.Price,
		"promotionPrice": product.PromotionPrice,
	}
	return bson.M{"$set": p}
}
func NewMongoProducts(products []*model.Product) []*mongoProduct {
	result := make([]*mongoProduct, len(products))
	for i, product := range products {
		result[i] = newMongoProduct(product)
	}
	return result
}

func (p *mongoProduct) ToProduct() *model.Product {
	return &model.Product{ID: model.ProductId(p.ProductID), Name: p.Name, Brand: p.Brand, Description: p.Description, Price: p.Price, PromotionPrice: p.PromotionPrice}
}
