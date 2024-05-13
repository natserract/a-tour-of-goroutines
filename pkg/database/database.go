package database

import (
	"context"
	"fmt"
	"goroutines/config"
	"goroutines/pkg/env"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

type DB struct {
	*pgxpool.Pool
	QueryBuilder *squirrel.StatementBuilderType
	url          string
}

func New(ctx context.Context, config *config.DB) (*DB, error) {
	pgUrl := `postgres://%s:%s@%s:%d/%s?%s`
	pgUrl = fmt.Sprintf(pgUrl,
		config.Username,
		config.Pass,
		config.Host,
		config.Port,
		config.Name,
		config.Params,
	)

	conf, err := pgxpool.ParseConfig(pgUrl)
	if err != nil {
		return nil, err
	}

	// Only show on development mode
	if !env.IsProduction() {
		conf.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger: &PGXStdLogger{
				slog.Default(),
			},
			LogLevel: tracelog.LogLevelInfo,
		}
	}

	// pgxpool default max number of connections is the number of CPUs on your machine returned by runtime.NumCPU().
	// This number is very conservative, and you might be able to improve performance for highly concurrent applications
	// by increasing it.
	// conf.MaxConns = runtime.NumCPU() * 5
	pool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("pgx connection error: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	return &DB{
		pool,
		&psql,
		pgUrl,
	}, nil
}

// ErrorCode returns the error code of the given error
func (db *DB) ErrorCode(err error) string {
	pgErr := err.(*pgconn.PgError)
	return pgErr.Code
}

// Close closes the database connection
func (db *DB) Close() {
	db.Pool.Close()
}
