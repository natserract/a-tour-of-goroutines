package controller

import (
	"errors"
	"goroutines/internal/product/errs"
	"goroutines/internal/product/request"
	"goroutines/internal/product/response"
	"goroutines/internal/product/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductController interface {
	CreateProduct(ctx *gin.Context)
	CreateProductGoroutines(ctx *gin.Context)
	CreateProductTx(ctx *gin.Context)
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
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if reqBody.ValidateProductCreate() != nil {
		validateErr := reqBody.ValidateProductCreate()

		ctx.AbortWithError(http.StatusBadRequest, validateErr)
		return
	}

	productCreated, err := c.svc.CreateProduct(&reqBody)
	if err != nil {
		switch {
		case errors.Is(err, errs.ProductErrsCategoryNotFound):
			ctx.AbortWithError(http.StatusBadRequest, err)
			break
		default:
			ctx.AbortWithError(http.StatusInternalServerError, err)
			break
		}

		return
	}

	productCreatedMappedResult := response.ProductToCreateResponse(productCreated)
	ctx.JSON(http.StatusCreated, productCreatedMappedResult)
}

func (c *productController) CreateProductGoroutines(ctx *gin.Context) {
	var reqBody request.ProductCreateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if reqBody.ValidateProductCreate() != nil {
		validateErr := reqBody.ValidateProductCreate()

		ctx.AbortWithError(http.StatusBadRequest, validateErr)
		return
	}

	productCreated := <-c.svc.CreateProductGoroutines(&reqBody)
	if productCreated.Error != nil {
		switch {
		case errors.Is(productCreated.Error, errs.ProductErrsCategoryNotFound):
			ctx.AbortWithError(http.StatusBadRequest, productCreated.Error)
			break
		default:
			ctx.AbortWithError(http.StatusInternalServerError, productCreated.Error)
			break
		}

		return
	}

	productCreatedMappedResult := response.ProductToCreateResponse(productCreated.Result)
	ctx.JSON(http.StatusCreated, productCreatedMappedResult)
}

func (c *productController) CreateProductTx(ctx *gin.Context) {
	var reqBody request.ProductCreateRequest
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if reqBody.ValidateProductCreate() != nil {
		validateErr := reqBody.ValidateProductCreate()

		ctx.AbortWithError(http.StatusBadRequest, validateErr)
		return
	}

	productCreated, err := c.svc.CreateProductTx(&reqBody)
	if err != nil {
		switch {
		case errors.Is(err, errs.ProductErrsCategoryNotFound):
			ctx.AbortWithError(http.StatusBadRequest, err)
			break
		default:
			ctx.AbortWithError(http.StatusInternalServerError, err)
			break
		}

		return
	}

	productCreatedMappedResult := response.ProductToCreateResponse(productCreated)
	ctx.JSON(http.StatusCreated, productCreatedMappedResult)
}
