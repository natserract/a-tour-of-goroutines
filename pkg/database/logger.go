package database

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/tracelog"
)

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
