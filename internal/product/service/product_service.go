package service

import (
	"context"
	categoryRepository "goroutines/internal/category/repository"
	"goroutines/internal/product"
	"goroutines/internal/product/errs"
	"goroutines/internal/product/repository"
	"goroutines/internal/product/request"
	"goroutines/pkg/database"

	"github.com/jackc/pgx/v5"
)

type ProductService interface {
	CreateProduct(p *request.ProductCreateRequest) (*product.Product, error)
}

type ProductDependency struct {
	Product  repository.ProductRepository
	Category categoryRepository.CategoryRepository
}

type productService struct {
	db   *database.DB
	repo *ProductDependency
	ctx  context.Context
}

func NewProductService(
	db *database.DB,
	repo *ProductDependency,
	ctx context.Context,
) ProductService {
	return &productService{
		db:   db,
		repo: repo,
		ctx:  ctx,
	}
}

func (svc *productService) CreateProduct(p *request.ProductCreateRequest) (*product.Product, error) {
	repo := svc.repo

	var result *product.Product
	if err := svc.db.BeginTransaction(svc.ctx, func(tx pgx.Tx, ctx context.Context) error {
		categoryFound, err := repo.Category.GetReferenceByName(ctx, p.Category)
		if err != nil {
			return errs.ProductErrorCategoryNotFound
		}

		product := &product.Product{
			Name:        p.Name,
			Sku:         p.Sku,
			Category:    categoryFound.Name,
			ImageUrl:    p.ImageUrl,
			Notes:       p.Notes,
			Price:       p.Price,
			Stock:       *p.Stock,
			Location:    p.Location,
			IsAvailable: *p.IsAvailable,
		}
		productPersisted, err := repo.Product.Persist(ctx, product, tx)
		if err != nil {
			return err
		}

		result = productPersisted
		return nil
	}); err != nil {
		return nil, err
	}

	return result, nil
}
