package router

import (
	"context"
	v1 "goroutines/router/v1"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRouter(ctx context.Context, pool *pgxpool.Pool, router *gin.Engine) {
	v1Route := v1.NewV1Router(ctx, pool)
	v1Route.Load(router)
}
