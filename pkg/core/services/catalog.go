package services

import (
	"context"
	"sync"

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
	err := model.ValidateProductIds(ids)
	if err != nil {
		return nil, err
	}
	return s.repo.GetProductByIds(ctx, ids)
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

func (s *catalogImportService) store(ctx context.Context, product *model.Product) error {
	err := model.ValidateProduct(product)
	if err != nil {
		return err
	}
	err = s.repo.Insert(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

func (s *catalogImportService) Store(ctx context.Context, products <-chan *model.Product) <-chan *model.Product {
	stream := make(chan *model.Product)
	go func() {
		var wg sync.WaitGroup
		for product := range products {
			wg.Add(1)
			go func(p *model.Product) {
				defer wg.Done()
				err := s.store(ctx, p)
				if err != nil {
					log.Errorf("error while storing product: %s", err)
				} else {
					stream <- p
				}
			}(product)
		}
		wg.Wait()
		close(stream)
	}()
	return stream
}

type ProductsSourceDataProvider interface {
	Provide(ctx context.Context) <-chan *model.Product
}

type productsSourceDataProvider struct {
}

func NewProductsSourceDataProvider() *productsSourceDataProvider {
	return &productsSourceDataProvider{}
}

func (s *productsSourceDataProvider) Provide(ctx context.Context) <-chan *model.Product {
	stream := make(chan *model.Product)
	go func() {
		for _, p := range []*model.Product{model.NewProduct(1, "xD", "c", "dsa", 23.23), model.NewPromotionalProduct(1, "xD2", "c2", "22", 23.23, 1.2)} {
			stream <- p
		}
		close(stream)
	}()
	return stream
}
