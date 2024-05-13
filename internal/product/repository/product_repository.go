package repository

import (
	"context"
	"goroutines/internal/product"
	"goroutines/pkg/database"

	"github.com/jackc/pgx/v5"
)

type ProductRepository interface {
	Persist(ctx context.Context, p *product.Product, tx pgx.Tx) (*product.Product, error)
}

type productRepository struct {
	db *database.DB
}

func NewProductRepository(db *database.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

// CreateProduct creates a new product record in the database
func (pr *productRepository) Persist(ctx context.Context, p *product.Product, tx pgx.Tx) (*product.Product, error) {
	query := pr.db.QueryBuilder.Insert("products").
		Columns("id", "name", "sku", "category", "image_url", "notes", "price", "stock", "location", "is_available", "created_at").
		Values(
			p.Id,
			p.Name,
			p.Sku,
			p.Category,
			p.ImageUrl,
			p.Notes,
			p.Price,
			p.Stock,
			p.Location,
			p.IsAvailable,
			p.CreatedAt,
		).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pr.db.QueryRow(ctx, sql, args...).Scan(
		&p.Id,
		&p.Name,
		&p.Sku,
		&p.Category,
		&p.ImageUrl,
		&p.Notes,
		&p.Price,
		&p.Stock,
		&p.Location,
		&p.IsAvailable,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errCode := pr.db.ErrorCode(err); errCode == "23505" {
			return nil, err
		}
		return nil, err
	}

	return p, nil
}
