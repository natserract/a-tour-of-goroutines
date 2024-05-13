package v1

import (
	"context"
	"goroutines/internal/product/controller"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRouter struct {
	Controller controller.ProductController
}

func NewProductRouter(ctx context.Context, pool *pgxpool.Pool) *ProductRouter {
	return &ProductRouter{
		Controller: controller.NewProductController(),
	}
}
