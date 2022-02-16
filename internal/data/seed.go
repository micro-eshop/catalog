package data

import (
	"context"
	"encoding/csv"
	"math/rand"
	"os"
	"strconv"

	"github.com/micro-eshop/catalog/pkg/core/model"
	log "github.com/sirupsen/logrus"
)

type productsSourceDataProvider struct {
	path string
}

func NewProductsSourceDataProvider(path string) *productsSourceDataProvider {
	return &productsSourceDataProvider{path}
}

func randomPrice(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func parseInt(idsStr string) model.ProductId {
	res, _ := strconv.Atoi(idsStr)
	return model.ProductId(res)
}

func parseFloat(idsStr string) float64 {
	res, _ := strconv.ParseFloat(idsStr, 64)
	return res
}

func (s *productsSourceDataProvider) Provide(ctx context.Context) <-chan *model.Product {
	stream := make(chan *model.Product)

	go func() {
		csvFile, err := os.Open(s.path)
		if err != nil {
			log.WithError(err).Fatalln("error while opening csv file")
		}
		log.Infoln("Successfully Opened CSV file")
		defer csvFile.Close()
		csvLines, err := csv.NewReader(csvFile).ReadAll()
		if err != nil {
			log.WithError(err).Fatalln("error while reading csv file")
		}
		for i, p := range csvLines {
			id := parseInt(p[0])
			price := parseFloat(p[5])

			if i%2 == 0 {
				stream <- model.NewProduct(id, p[2], p[3], p[4], price)
			} else {
				stream <- model.NewPromotionalProduct(id, p[2], p[3], p[4], price, randomPrice(1, 10))
			}
		}
		close(stream)
	}()
	return stream
}
