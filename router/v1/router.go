package v1

import (
	"context"
	"goroutines/pkg/database"

	"github.com/gin-gonic/gin"
)

type V1Router interface {
	Load(r *gin.Engine)
}

type v1Router struct {
	Product *ProductRouter
}

func NewV1Router(ctx context.Context, db *database.DB) *v1Router {
	return &v1Router{
		Product: NewProductRouter(ctx, db),
	}
}

func (v *v1Router) Load(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		// Product api endpoint
		product := v1.Group("/product")
		product.POST("/", v.Product.Controller.CreateProductGoroutines)
	}
}
