package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/micro-eshop/catalog/pkg/core/model"
	"github.com/micro-eshop/catalog/pkg/core/repositories"
	log "github.com/sirupsen/logrus"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type postgresProduct struct {
	ProductID      int             `pg:"id"`
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
	db *sql.DB
}

func createSchema(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgresql", driver)
	if err != nil {
		return err
	}
	err = m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
	if err != nil && err.Error() == "no change" {
		return nil
	}
	return err
}

func NewPostgresClient(ctx context.Context, connectionString string) (*postgresClient, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
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
	sql, _, err := sq.Select("id", "brand", "name", "description", "price", "promotion_price").From("products").Where(sq.Eq{"id": int(id)}).ToSql()
	if err != nil {
		return nil, err
	}
	var dbProduct postgresProduct
	row := r.client.db.QueryRowContext(ctx, sql)
	row.Scan(&dbProduct)
	if err != nil {
		return nil, err
	}

	return dbProduct.toProduct(), nil
}

func (r *postgresCatalogRepository) GetProductByIds(ctx context.Context, ids []model.ProductId) ([]*model.Product, error) {
	return nil, nil
}

func (r *postgresCatalogRepository) Search(ctx context.Context, params repositories.ProductSearchParams) ([]*model.Product, error) {
	return nil, nil
}

func (r *postgresCatalogRepository) Insert(ctx context.Context, product *model.Product) error {
	dbProduct := newPostgresProduct(product)
	log.Println("Inserting product: ", dbProduct)
	query := sq.Insert("products").
		Columns("brand", "name", "description", "price", "promotion_price").
		Values(dbProduct.Brand, dbProduct.Name, dbProduct.Description, dbProduct.Price, dbProduct.PromotionPrice).
		RunWith(r.client.db).
		PlaceholderFormat(sq.Dollar)

	_, err := query.Exec()

	return err
}
