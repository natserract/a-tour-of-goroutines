package database

import (
	"context"
	"fmt"
	"goroutines/pkg/env"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
)

// NewPGXPool is a PostgreSQL connection pool for pgx.
//
// Usage:
// pgPool := database.NewPGXPool(context.Background(), "", &PGXStdLogger{Logger: slog.Default()}, tracelog.LogLevelInfo, tracer)
// defer pgPool.Close() // Close any remaining connections before shutting down your application.
//
// Instead of passing a configuration explicitly with a connString,
// you might use PG environment variables such as the following to configure the database:
// PGDATABASE, PGHOST, PGPORT, PGUSER, PGPASSWORD, PGCONNECT_TIMEOUT, etc.
// Reference: https://www.postgresql.org/docs/current/libpq-envars.html
func NewPGXPool(ctx context.Context, connString string, logger tracelog.Logger) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(connString) // Using environment variables instead of a connection string.
	if err != nil {
		return nil, err
	}

	// Only show on development mode
	if !env.IsProduction() {
		conf.ConnConfig.Tracer = &tracelog.TraceLog{
			Logger:   logger,
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

	fmt.Printf("Successfully connected to database %s %v", connString, "\n")
	return pool, nil

}

// PGXStdLogger prints pgx logs to the standard logger.
// os.Stderr by default.
type PGXStdLogger struct {
	Logger *slog.Logger
}

func (l *PGXStdLogger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	attrs := make([]slog.Attr, 0, len(data)+1)
	attrs = append(attrs, slog.String("pgx_level", level.String()))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}
	l.Logger.LogAttrs(ctx, slogLevel(level), msg, attrs...)
}

// slogLevel translates pgx log level to slog log level.
func slogLevel(level tracelog.LogLevel) slog.Level {
	switch level {
	case tracelog.LogLevelTrace, tracelog.LogLevelDebug:
		return slog.LevelDebug
	case tracelog.LogLevelInfo:
		return slog.LevelInfo
	case tracelog.LogLevelWarn:
		return slog.LevelWarn
	default:
		// If tracelog.LogLevelError, tracelog.LogLevelNone, or any other unknown level, use slog.LevelError.
		return slog.LevelError
	}
}
