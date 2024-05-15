package repository

import (
	"context"
	"errors"
	domain "goroutines/internal/category"
	"goroutines/pkg/database"
	"log/slog"

	sq "github.com/Masterminds/squirrel"

	"github.com/jackc/pgx/v5"
)

type CategoryRepository interface {
	GetReferenceByName(ctx context.Context, name string) (*domain.Category, error)
}

type categoryRepository struct {
	db *database.DB
}

func NewCategoryRepository(db *database.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (cr *categoryRepository) GetReferenceByName(ctx context.Context, name string) (*domain.Category, error) {
	var category domain.Category

	query := cr.db.QueryBuilder.Select("*").
		From("categories").
		Where(sq.Eq{"name": name}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cr.db.Query(ctx, sql, args...)
	if err == nil {
		category, err = pgx.CollectOneRow(rows, pgx.RowToStructByPos[domain.Category])
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if err != nil {
		slog.Error("cannot get category from database",
			slog.Any("name", name),
			slog.Any("error", err))
		return nil, errors.New("cannot get category from database")
	}

	return &category, nil
}
