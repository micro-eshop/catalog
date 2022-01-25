package mongo

import (
	"context"

	"github.com/micro-eshop/catalog/internal/core/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func ToInsertMongoDocument(product *mongoProduct) bson.M {
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

func ToInterfaceSlice(ss []*mongoProduct) []interface{} {
	iface := make([]interface{}, len(ss))
	for i := range ss {
		iface[i] = ss[i]
	}
	return iface
}

type mongoCatalogRepository struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoCatalogRepository(client *mongo.Client) *mongoCatalogRepository {
	db := client.Database(catalogDatabase)
	return &mongoCatalogRepository{client: client, db: db}
}

func (r *mongoCatalogRepository) GetProductById(ctx context.Context, id model.ProductId) (*model.Product, error) {
	col := r.db.Collection(productsCollectionName, &options.CollectionOptions{})
	filter := bson.M{"_id": id}
	var result mongoProduct
	err := col.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return result.ToProduct(), nil
}

func (r *mongoCatalogRepository) GetProductByIds(ctx context.Context, ids []model.ProductId) ([]*model.Product, error) {
	col := r.db.Collection(productsCollectionName, &options.CollectionOptions{})
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cur, err := col.Find(ctx, filter)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var results []*model.Product
	for cur.Next(ctx) {
		var result mongoProduct
		err = cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result.ToProduct())

	}
	return results, nil
}
