package service

import (
	"context"
	"goroutines/internal/product"
	"goroutines/internal/product/repository"
)

type ProductService interface {
	CreateProduct(ctx context.Context, p *product.Product) (*product.Product, error)
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{
		productRepository: productRepo,
	}
}

func (svc *productService) CreateProduct(ctx context.Context, p *product.Product) (*product.Product, error) {
	return nil, nil
}
