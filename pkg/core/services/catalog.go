package services

import (
	"context"
	"math/rand"
	"sync"

	"github.com/bxcodec/faker/v3"
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

func randomPrice(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func (s *productsSourceDataProvider) Provide(ctx context.Context) <-chan *model.Product {
	stream := make(chan *model.Product)

	go func() {
		for i := 0; i < 100; i++ {
			if i%2 == 0 {
				var f float64 = randomPrice(float64(i), 200)
				stream <- model.NewPromotionalProduct(model.ProductId(i+1), faker.Name(), faker.Word(), faker.Sentence(), randomPrice(200, 400), f)
			} else {
				stream <- model.NewProduct(model.ProductId(i+1), faker.Name(), faker.Word(), faker.Sentence(), randomPrice(200, 400))
			}
		}
		close(stream)
	}()
	return stream
}
