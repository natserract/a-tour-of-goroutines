package v1

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type V1Router interface {
	Load(r *gin.Engine)
}

type v1Router struct {
	Product *ProductRouter
}

func NewV1Router(ctx context.Context, pool *pgxpool.Pool) *v1Router {
	return &v1Router{
		Product: NewProductRouter(ctx, pool),
	}
}

func (v *v1Router) Load(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		// Product api endpoint
		product := v1.Group("/product")
		product.POST("/", v.Product.Controller.CreateProduct)
	}
}
