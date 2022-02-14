package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/micro-eshop/catalog/pkg/core/model"
	"github.com/micro-eshop/catalog/pkg/core/repositories"
	log "github.com/sirupsen/logrus"

	sq "github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

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

func newNullFloat64(s *float64) sql.NullFloat64 {
	if s == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: float64(*s),
		Valid:   true,
	}
}

func newPostgresProduct(product *model.Product) *postgresProduct {
	return &postgresProduct{ProductID: int(product.ID), Name: product.Name, Brand: product.Brand, Description: product.Description, Price: product.Price, PromotionPrice: newNullFloat64(product.PromotionPrice)}
}

func (p postgresProduct) toProduct() *model.Product {
	var price *float64
	if p.PromotionPrice.Valid {
		price = &p.PromotionPrice.Float64
	}
	return &model.Product{ID: model.ProductId(p.ProductID), Name: p.Name, Brand: p.Brand, Description: p.Description, Price: p.Price, PromotionPrice: price}
}

func mapProduct(scanner sq.RowScanner) (*postgresProduct, error) {
	var dbProduct postgresProduct
	err := scanner.Scan(&dbProduct.ProductID, &dbProduct.Brand, &dbProduct.Name, &dbProduct.Description, &dbProduct.Price, &dbProduct.PromotionPrice)
	if err != nil {
		return nil, err
	}
	return &dbProduct, nil
}

func mapIds(ids []model.ProductId) []int {
	result := make([]int, len(ids))
	for i, id := range ids {
		result[i] = int(id)
	}
	return result
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
	query := psql.Select("id", "brand", "name", "description", "price", "promotion_price").From("products").Where(sq.Eq{"id": int(id)})
	row := query.RunWith(r.client.db).QueryRowContext(ctx)
	product, err := mapProduct(row)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return product.toProduct(), nil
}

func (r *postgresCatalogRepository) GetProductByIds(ctx context.Context, ids []model.ProductId) ([]*model.Product, error) {
	query := psql.Select("id", "brand", "name", "description", "price", "promotion_price").From("products").Where(sq.Eq{"id": mapIds(ids)})
	rows, err := query.RunWith(r.client.db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	products := make([]*model.Product, 0)
	for rows.Next() {
		product, nerr := mapProduct(rows)
		if nerr != nil {
			err = multierror.Append(nerr)
		}
		products = append(products, product.toProduct())
	}

	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, nil
	}
	return products, nil
}

func (r *postgresCatalogRepository) Search(ctx context.Context, params repositories.ProductSearchParams) ([]*model.Product, error) {
	return nil, nil
}

func (r *postgresCatalogRepository) Insert(ctx context.Context, product *model.Product) error {
	dbProduct := newPostgresProduct(product)
	log.Println("Inserting product: ", dbProduct)
	query := psql.Insert("products").
		Columns("brand", "name", "description", "price", "promotion_price").
		Values(dbProduct.Brand, dbProduct.Name, dbProduct.Description, dbProduct.Price, dbProduct.PromotionPrice).
		RunWith(r.client.db)

	_, err := query.Exec()

	return err
}
