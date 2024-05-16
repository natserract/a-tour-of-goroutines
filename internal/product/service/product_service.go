package service

import (
	"context"
	categoryRepository "goroutines/internal/category/repository"
	"goroutines/internal/product"
	"goroutines/internal/product/errs"
	"goroutines/internal/product/repository"
	"goroutines/internal/product/request"
	"goroutines/pkg/database"
	"goroutines/util"
	"runtime"

	"github.com/jackc/pgx/v5"
)

type ProductService interface {
	CreateProduct(p *request.ProductCreateRequest) (*product.Product, error)
	CreateProductGoroutines(p *request.ProductCreateRequest) <-chan util.Result[*product.Product]
	CreateProductGoroutinesIncrease(p *request.ProductCreateRequest) <-chan util.Result[*product.Product]
	CreateProductTx(p *request.ProductCreateRequest) (*product.Product, error)
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

	categoryFound, err := repo.Category.GetReferenceByName(svc.ctx, p.Category)
	if err != nil {
		return nil, errs.ProductErrsCategoryNotFound
	}

	model := &product.Product{
		Name:        p.Name,
		Sku:         p.Sku,
		Category:    categoryFound.Name,
		ImageUrl:    p.ImageUrl,
		Notes:       p.Notes,
		Price:       p.Price,
		Stock:       *p.Stock,
		Location:    p.Location,
		IsAvailable: p.IsAvailable,
	}
	productPersisted, err := repo.Product.Persist(svc.ctx, model)
	if err != nil {
		return nil, err
	}

	return productPersisted, nil
}

func (svc *productService) CreateProductGoroutines(p *request.ProductCreateRequest) <-chan util.Result[*product.Product] {
	repo := svc.repo

	result := make(chan util.Result[*product.Product])
	go func() {
		categoryFound, err := repo.Category.GetReferenceByName(svc.ctx, p.Category)
		if err != nil {
			result <- util.Result[*product.Product]{
				Error: errs.ProductErrsCategoryNotFound,
			}
			return
		}

		model := &product.Product{
			Name:        p.Name,
			Sku:         p.Sku,
			Category:    categoryFound.Name,
			ImageUrl:    p.ImageUrl,
			Notes:       p.Notes,
			Price:       p.Price,
			Stock:       *p.Stock,
			Location:    p.Location,
			IsAvailable: p.IsAvailable,
		}
		productPersisted, err := repo.Product.Persist(svc.ctx, model)
		if err != nil {
			result <- util.Result[*product.Product]{
				Error: err,
			}
			return
		}

		result <- util.Result[*product.Product]{
			Result: productPersisted,
		}
		close(result)
	}()

	return result
}

func (svc *productService) CreateProductGoroutinesIncrease(p *request.ProductCreateRequest) <-chan util.Result[*product.Product] {
	repo := svc.repo

	worker := runtime.NumCPU()
	result := make(chan util.Result[*product.Product], worker)

	task := func() {
		defer close(result)

		categoryFound, err := repo.Category.GetReferenceByName(svc.ctx, p.Category)
		if err != nil {
			result <- util.Result[*product.Product]{
				Error: errs.ProductErrsCategoryNotFound,
			}
			return
		}

		model := &product.Product{
			Name:        p.Name,
			Sku:         p.Sku,
			Category:    categoryFound.Name,
			ImageUrl:    p.ImageUrl,
			Notes:       p.Notes,
			Price:       p.Price,
			Stock:       *p.Stock,
			Location:    p.Location,
			IsAvailable: p.IsAvailable,
		}
		productPersisted, err := repo.Product.Persist(svc.ctx, model)
		if err != nil {
			result <- util.Result[*product.Product]{
				Error: err,
			}
			return
		}

		result <- util.Result[*product.Product]{
			Result: productPersisted,
		}
	}
	go task()

	return result
}

func (svc *productService) CreateProductTx(p *request.ProductCreateRequest) (*product.Product, error) {
	repo := svc.repo

	var result *product.Product
	if err := svc.db.BeginTransaction(svc.ctx, func(tx pgx.Tx, ctx context.Context) error {
		categoryFound, err := repo.Category.GetReferenceByName(ctx, p.Category)
		if err != nil {
			return errs.ProductErrsCategoryNotFound
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
			IsAvailable: p.IsAvailable,
		}
		productPersisted, err := repo.Product.PersistTx(ctx, product, tx)
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
