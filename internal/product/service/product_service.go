package service

import (
	"context"
	"fmt"
	categoryRepository "goroutines/internal/category/repository"
	"goroutines/internal/product"
	"goroutines/internal/product/repository"
	"goroutines/internal/product/request"
)

type ProductService interface {
	CreateProduct(ctx context.Context, p *request.ProductCreateRequest) (*product.Product, error)
}

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo categoryRepository.CategoryRepository
}

func NewProductService(
	productRepo repository.ProductRepository,
	categoryRepo categoryRepository.CategoryRepository,
) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (svc *productService) validate(r *request.ProductCreateRequest) error {
	return nil
}

func (svc *productService) CreateProduct(ctx context.Context, p *request.ProductCreateRequest) (*product.Product, error) {
	category, err := svc.categoryRepo.GetReferenceByName(ctx, p.Category)
	if err != nil {
		return nil, err
	}

	fmt.Println("category", category)
	return nil, nil
}
