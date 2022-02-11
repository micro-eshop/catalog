package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/micro-eshop/catalog/pkg/core/model"
	"github.com/micro-eshop/catalog/pkg/core/usecase"
)

type CatalogHandler struct {
	GetProductByIdUseCase *usecase.GetProductByIdUseCase
}

func NewCatalogHandler(getProductByIdUseCase *usecase.GetProductByIdUseCase) *CatalogHandler {
	return &CatalogHandler{
		GetProductByIdUseCase: getProductByIdUseCase,
	}
}

func (handler *CatalogHandler) GetProductById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "id is not a number",
		})
		return
	}
	productId := model.ProductId(id)
	product, err := handler.GetProductByIdUseCase.Execute(c.Request.Context(), productId)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if product == nil {
		c.JSON(404, gin.H{
			"message": "product not found",
		})
		return
	}
	c.JSON(200, product)
}
