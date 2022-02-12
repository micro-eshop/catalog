package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/micro-eshop/catalog/pkg/core/model"
	"github.com/micro-eshop/catalog/pkg/core/usecase"
)

func parseInt(idsStr []string) []model.ProductId {
	res := make([]model.ProductId, len(idsStr))
	for i, v := range idsStr {
		id, _ := strconv.Atoi(v)
		res[i] = model.ProductId(id)
	}
	return res
}

type CatalogHandler struct {
	getProductByIdUseCase  *usecase.GetProductByIdUseCase
	getProductByIdsUseCase *usecase.GetProductByIdsUseCase
}

func NewCatalogHandler(getProductByIdUseCase *usecase.GetProductByIdUseCase, getProductByIdsUseCase *usecase.GetProductByIdsUseCase) *CatalogHandler {
	return &CatalogHandler{
		getProductByIdUseCase:  getProductByIdUseCase,
		getProductByIdsUseCase: getProductByIdsUseCase,
	}
}

func (handler *CatalogHandler) getProductById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "id is not a number",
		})
		return
	}
	productId := model.ProductId(id)
	product, err := handler.getProductByIdUseCase.Execute(c.Request.Context(), productId)
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

func (handler *CatalogHandler) getProductByIds(c *gin.Context) {
	idsStr := c.QueryArray("ids")

	if len(idsStr) == 0 {
		c.Status(204)
		return
	}

	ids := parseInt(idsStr)

	result, err := handler.getProductByIdsUseCase.Execute(c.Request.Context(), ids)

	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "unknown error")
		return
	}

	if result == nil {
		c.Status(404)
		return
	}

	c.JSON(200, result)
}

func (h *CatalogHandler) Setup(r gin.IRouter) {
	r.Group("/catalog").GET("/products/:id", h.getProductById).GET("/products", h.getProductByIds)
}
