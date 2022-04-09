package services

import (
	"context"

	"github.com/micro-eshop/catalog/pkg/core/model"
	"github.com/micro-eshop/catalog/pkg/core/repositories"
	log "github.com/sirupsen/logrus"
)

type CatalogService interface {
	GetProductById(ctx context.Context, id model.ProductId) (*model.Product, error)
	GetProductByIds(ctx context.Context, ids []model.ProductId) ([]*model.Product, error)
	Search(ctx context.Context, params repositories.ProductSearchParams) ([]*model.Product, error)
}

type catalogService struct {
	repo repositories.CatalogReader
}

func NewCatalogService(repo repositories.CatalogReader) *catalogService {
	return &catalogService{
		repo: repo,
	}
}

func (s *catalogService) GetProductById(ctx context.Context, id model.ProductId) (*model.Product, error) {
	err := model.ValidateProductId(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetProductById(ctx, id)
}

func (s *catalogService) GetProductByIds(ctx context.Context, ids []model.ProductId) ([]*model.Product, error) {
	if len(ids) == 0 {
		return make([]*model.Product, 0), nil
	}

	err := model.ValidateProductIds(ids)
	if err != nil {
		return nil, err
	}
	return s.repo.GetProductByIds(ctx, ids...)
}

func (s *catalogService) Search(ctx context.Context, params repositories.ProductSearchParams) ([]*model.Product, error) {
	return s.repo.Search(ctx, params)
}

type CatalogImportService interface {
	Store(ctx context.Context, products <-chan *model.Product) <-chan *model.Product
}

type catalogImportService struct {
	repo repositories.CatalogWriter
}

func NewCatalogImportService(repo repositories.CatalogRepository) *catalogImportService {
	return &catalogImportService{
		repo: repo,
	}
}

func filter(ctx context.Context, products <-chan *model.Product, predicate func(ctx context.Context, product *model.Product) bool) <-chan *model.Product {
	filteredProducts := make(chan *model.Product, 100)
	go func() {
		for product := range products {
			if predicate(ctx, product) {
				filteredProducts <- product
			}
		}
		close(filteredProducts)
	}()
	return filteredProducts
}

func pipe(ctx context.Context, products <-chan *model.Product, f func(ctx context.Context, product *model.Product, stream chan<- *model.Product)) <-chan *model.Product {
	out := make(chan *model.Product, 10)
	go func() {
		for product := range products {
			f(ctx, product, out)
		}
		close(out)
	}()
	return out
}

func (s *catalogImportService) Store(ctx context.Context, products <-chan *model.Product) <-chan *model.Product {
	validProducts := filter(ctx, products, func(ctx context.Context, product *model.Product) bool {
		err := model.ValidateProduct(product)
		if err != nil {
			log.WithError(err).WithField("ProductId", product.ID).Errorln("Product is not valid")
			return false
		}
		return true
	})

	return pipe(ctx, validProducts, func(ctx context.Context, product *model.Product, out chan<- *model.Product) {
		err := s.repo.Insert(ctx, product)
		if err != nil {
			log.Errorf("Error storing product: %v", err)
		} else {
			out <- product
		}
	})
}

type ProductsSourceDataProvider interface {
	Provide(ctx context.Context) <-chan *model.Product
}
