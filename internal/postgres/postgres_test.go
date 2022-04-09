package postgres

import (
	"context"
	"testing"

	"github.com/dominikus1993/integrationtestcontainers-go"
	"github.com/micro-eshop/catalog/pkg/core/model"
	"github.com/stretchr/testify/assert"
)

func TestNewNullFloat64WhenNil(t *testing.T) {
	var id *float64 = nil
	subject := newNullFloat64(id)
	assert.NotNil(t, subject)
	assert.False(t, subject.Valid)
	assert.Zero(t, subject.Float64)
}

func TestNewNullFloat64WhenNotNil(t *testing.T) {
	val := float64(1.0)
	var id *float64 = &val
	subject := newNullFloat64(id)
	assert.NotNil(t, subject)
	assert.True(t, subject.Valid)
	assert.Equal(t, 1.0, subject.Float64)
}

func TestGetPromotionPriceWhenNil(t *testing.T) {
	val := newPostgresProduct(&model.Product{ID: model.ProductId(1), Name: "name", Brand: "brand", Description: "description", Price: 1.0, PromotionPrice: nil})
	subject := val.getPromotionPrice()
	assert.Nil(t, subject)
}

func TestGetPromotionPriceWhenNotNil(t *testing.T) {
	price := float64(1.0)
	val := newPostgresProduct(&model.Product{ID: model.ProductId(1), Name: "name", Brand: "brand", Description: "description", Price: 1.0, PromotionPrice: &price})
	subject := val.getPromotionPrice()
	assert.NotNil(t, subject)
	assert.Equal(t, price, *subject)
}
func TestGetProductById(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	postgres, err := integrationtestcontainers.StartPostgreSqlContainer(ctx, integrationtestcontainers.DefaultPostgresContainerConfiguration)
	if err != nil {
		t.Fatal(err)
	}
	defer postgres.Terminate(ctx)
	db, err := NewPostgresClient(ctx, postgres.ConnectionString)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(ctx)
	repository := NewPostgresCatalogRepository(db)
	t.Run("when product does not exist", func(t *testing.T) {
		product, _ := repository.GetProductById(ctx, model.ProductId(1))
		assert.Nil(t, product)
	})

	t.Run("when product does exist", func(t *testing.T) {
		productId := model.ProductId(2)
		product := model.NewProduct(productId, "name", "brand", "description", 1.0)
		err := repository.Insert(ctx, product)
		if err != nil {
			t.Error(err)
		}
		dbproduct, err := repository.GetProductById(ctx, product.ID)
		assert.Nil(t, err)
		assert.NotNil(t, dbproduct)
		assert.Equal(t, product.ID, dbproduct.ID)
		assert.Equal(t, product.Name, dbproduct.Name)
		assert.Equal(t, product.Brand, dbproduct.Brand)
		assert.Equal(t, product.Description, dbproduct.Description)
		assert.Equal(t, product.Price, dbproduct.Price)
	})
}

func TestGetProductByIds(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	postgres, err := integrationtestcontainers.StartPostgreSqlContainer(ctx, integrationtestcontainers.DefaultPostgresContainerConfiguration)
	if err != nil {
		t.Fatal(err)
	}
	defer postgres.Terminate(ctx)
	db, err := NewPostgresClient(ctx, postgres.ConnectionString)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(ctx)
	repository := NewPostgresCatalogRepository(db)
	t.Run("when product does not exist", func(t *testing.T) {
		product, _ := repository.GetProductByIds(ctx, model.ProductId(1), model.ProductId(2))
		assert.Nil(t, product)
	})

	t.Run("when product does exist", func(t *testing.T) {
		productId := model.ProductId(2)
		productId2 := model.ProductId(3)
		product := model.NewProduct(productId, "name", "brand", "description", 1.0)
		product2 := model.NewProduct(productId2, "name", "brand", "description", 1.0)
		err := repository.Insert(ctx, product)
		if err != nil {
			t.Error(err)
		}
		err = repository.Insert(ctx, product2)
		if err != nil {
			t.Error(err)
		}
		products, err := repository.GetProductByIds(ctx, product.ID, product2.ID)
		assert.Nil(t, err)
		assert.NotNil(t, products)
		assert.Equal(t, 2, len(products))
		assert.Contains(t, products, product)
		assert.Contains(t, products, product2)
	})
}
