package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BeginTransaction(ctx context.Context, pool *pgxpool.Pool, f func(tx pgx.Tx, ctx context.Context) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Begin %w", err)
	}

	if err := f(tx, ctx); err != nil {
		_ = tx.Rollback(ctx)

		return fmt.Errorf("f %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("Commit %w", err)
	}

	return nil
}
