package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController interface {
	CreateProduct(ctx *gin.Context)
}

type productController struct{}

func NewProductController() ProductController {
	return &productController{}
}

func (c *productController) CreateProduct(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "OK")
}
