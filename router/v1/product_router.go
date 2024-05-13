package v1

import (
	"context"
	categoryRepository "goroutines/internal/category/repository"
	"goroutines/internal/product/controller"
	"goroutines/internal/product/repository"
	"goroutines/internal/product/service"
	"goroutines/pkg/database"
)

type ProductRouter struct {
	Controller controller.ProductController
}

func NewProductRouter(ctx context.Context, db *database.DB) *ProductRouter {
	productRepo := repository.NewProductRepository(db)
	categoryRepo := categoryRepository.NewCategoryRepository(db)
	productService := service.NewProductService(productRepo, categoryRepo)

	return &ProductRouter{
		Controller: controller.NewProductController(productService),
	}
}
