package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/micro-eshop/catalog/internal/core/model"
	"github.com/micro-eshop/catalog/internal/core/repositories"
)

type postgresProduct struct {
	ProductID      int             `pg:"id,pk"`
	Name           string          `pg:"name"`
	Brand          string          `pg:"brand"`
	Description    string          `pg:"description"`
	Price          float64         `pg:"price"`
	PromotionPrice sql.NullFloat64 `pg:"promotion_price"`
}

func (u postgresProduct) String() string {
	return fmt.Sprintf("Product<%d, %s, %s, %s, %f, %f)>", u.ProductID, u.Name, u.Brand, u.Description, u.Price, u.PromotionPrice.Float64)
}

func newPostgresProduct(product *model.Product) *postgresProduct {
	var promotionalPrice sql.NullFloat64
	if product.PromotionPrice != nil {
		promotionalPrice.Float64 = *product.PromotionPrice
	}
	return &postgresProduct{ProductID: int(product.ID), Name: product.Name, Brand: product.Brand, Description: product.Description, Price: product.Price, PromotionPrice: promotionalPrice}
}

func (p postgresProduct) toProduct() *model.Product {
	var price *float64
	if p.PromotionPrice.Valid {
		price = &p.PromotionPrice.Float64
	}
	return &model.Product{ID: model.ProductId(p.ProductID), Name: p.Name, Brand: p.Brand, Description: p.Description, Price: p.Price, PromotionPrice: price}
}

type postgresClient struct {
	db *pg.DB
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*postgresProduct)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func NewPostgresClient(ctx context.Context, connectionString string) (*postgresClient, error) {
	opt, err := pg.ParseURL(connectionString)
	if err != nil {
		return nil, err
	}
	db := pg.Connect(opt)
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	err = createSchema(db)
	if err != nil {
		return nil, err
	}
	return &postgresClient{db: db}, nil
}

func (c *postgresClient) Close(_ context.Context) {
	if err := c.db.Close(); err != nil {
		panic(err)
	}
}

type postgresCatalogRepository struct {
	client *postgresClient
}

func NewPostgresCatalogRepository(postgresClient *postgresClient) *postgresCatalogRepository {
	return &postgresCatalogRepository{client: postgresClient}
}

func (r *postgresCatalogRepository) GetProductById(ctx context.Context, id model.ProductId) (*model.Product, error) {
	product := &postgresProduct{ProductID: int(id)}
	err := r.client.db.Model(product).WherePK().Limit(1).Select()

	if err != nil {
		return nil, err
	}

	return product.toProduct(), nil
}

func (r *postgresCatalogRepository) GetProductByIds(ctx context.Context, ids []model.ProductId) ([]*model.Product, error) {
	return nil, nil
}

func (r *postgresCatalogRepository) Search(ctx context.Context, params repositories.ProductSearchParams) ([]*model.Product, error) {
	return nil, nil
}

func (r *postgresCatalogRepository) Insert(ctx context.Context, product *model.Product) error {
	dbProduct := newPostgresProduct(product)
	_, err := r.client.db.Model(dbProduct).OnConflict("(id) DO UPDATE").Insert()
	return err
}
