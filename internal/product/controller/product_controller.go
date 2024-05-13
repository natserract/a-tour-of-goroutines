package controller

import (
	"goroutines/internal/product/request"
	"goroutines/internal/product/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController interface {
	CreateProduct(ctx *gin.Context)
}

type productController struct {
	svc service.ProductService
}

func NewProductController(svc service.ProductService) ProductController {
	return &productController{svc}
}

func (c *productController) CreateProduct(ctx *gin.Context) {
	var reqBody request.ProductCreateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := c.svc.CreateProduct(ctx.Copy(), &reqBody)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "OK")
}
