package repository

import (
	"context"
	"errors"
	domain "goroutines/internal/product"
	"goroutines/pkg/database"
	"log/slog"

	sq "github.com/Masterminds/squirrel"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type ProductRepository interface {
	GetReferenceById(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	Persist(ctx context.Context, p *domain.Product) (*domain.Product, error)
	PersistTx(ctx context.Context, p *domain.Product, parentTx pgx.Tx) (*domain.Product, error)
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
func (pr *productRepository) Persist(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	db := pr.db
	query := db.QueryBuilder.Insert("products").
		Columns("name", "sku", "category", "image_url", "notes", "price", "stock", "location", "is_available", "created_at").
		Values(
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
		Suffix("RETURNING id, name, sku, category, image_url, notes, price, stock, location, is_available, created_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = db.Pool.QueryRow(ctx, sql, args...).Scan(
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
	)
	if err != nil {
		if sqlErr := pr.db.ErrorCode(err); sqlErr != nil {
			return nil, sqlErr
		}

		slog.Error("Cannot persist product on database", slog.Any("error", err))
		return nil, err
	}

	return p, nil
}

func (pr *productRepository) PersistTx(ctx context.Context, p *domain.Product, parentTx pgx.Tx) (*domain.Product, error) {
	db := pr.db
	query := db.QueryBuilder.Insert("products").
		Columns("name", "sku", "category", "image_url", "notes", "price", "stock", "location", "is_available", "created_at").
		Values(
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
		Suffix("RETURNING id, name, sku, category, image_url, notes, price, stock, location, is_available, created_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = parentTx.QueryRow(ctx, sql, args...).Scan(
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
	)
	if err != nil {
		if sqlErr := pr.db.ErrorCode(err); sqlErr != nil {
			return nil, sqlErr
		}

		slog.Error("Cannot persist product on database", slog.Any("error", err))
		return nil, err
	}

	return p, nil
}

func (pr *productRepository) GetReferenceById(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var p domain.Product

	query := pr.db.QueryBuilder.Select(`
		id, name, sku, category, image_url, notes, price, stock,
		location, is_available, created_at, updated_at, deleted_at
	`).
		From("products").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pr.db.Query(ctx, sql, args...)
	if err == nil {
		p, err = pgx.CollectOneRow(rows, pgx.RowToStructByPos[domain.Product])
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		slog.Error("cannot get product from database",
			slog.Any("id", id),
			slog.Any("error", err))
		return nil, errors.New("cannot get product from database")
	}

	return &p, nil
}
