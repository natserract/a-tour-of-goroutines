package router

import (
	"context"
	"goroutines/pkg/database"
	v1 "goroutines/router/v1"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(ctx context.Context, db *database.DB, router *gin.Engine) {
	v1Route := v1.NewV1Router(ctx, db)
	v1Route.Load(router)
}
