package repositories

import (
	"context"

	"github.com/micro-eshop/catalog/internal/core/model"
	"github.com/micro-eshop/catalog/internal/core/repositories"
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

type mongoCatalogRepository struct {
	client *MongoClient
	db     *mongo.Database
}

func NewMongoCatalogRepository(client *MongoClient) *mongoCatalogRepository {
	db := client.mongo.Database(catalogDatabase)
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

func (r *mongoCatalogRepository) Search(ctx context.Context, params repositories.ProductSearchParams) ([]*model.Product, error) {
	col := r.db.Collection(productsCollectionName, &options.CollectionOptions{})
	filter := bson.M{}
	if params.Name != "" {
		filter["name"] = bson.M{"$regex": params.Name}
	}
	if params.Brand != "" {
		filter["brand"] = bson.M{"$regex": params.Brand}
	}
	if params.PriceFrom != 0 {
		filter["price"] = bson.M{"$gte": params.PriceFrom}
	}
	if params.PriceTo != 0 {
		filter["price"] = bson.M{"$lte": params.PriceTo}
	}
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

func (r *mongoCatalogRepository) Insert(ctx context.Context, product *model.Product) error {
	col := r.db.Collection(productsCollectionName)
	opt := options.Update().SetUpsert(true)
	dbProduct := newMongoProduct(product)
	filter := bson.M{"_id": product.ID}
	_, err := col.UpdateOne(ctx, filter, toInsertMongoDocument(dbProduct), opt)
	return err
}
