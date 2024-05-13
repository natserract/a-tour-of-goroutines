package main

import (
	"context"
	"errors"
	"fmt"
	"goroutines/config"
	"goroutines/pkg/database"
	"goroutines/pkg/env"
	"log/slog"
	"net/http"
	"os"

	routes "goroutines/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load .env %v\n", err)
		os.Exit(1)
	}

	cfg := config.GetConfig()

	// Shared ctx
	ctx := context.Background()

	// Connect to the database
	pgUrl := `postgres://%s:%s@%s:%d/%s?%s`
	pgUrl = fmt.Sprintf(pgUrl,
		cfg.DBUsername,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBParams,
	)
	pool, err := database.NewPGXPool(ctx, pgUrl, &database.PGXStdLogger{
		Logger: slog.Default(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Disable debug mode in production
	if env.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Check reachability
	if _, err = pool.Exec(ctx, `SELECT 1`); err != nil {
		errMsg := fmt.Errorf("pool.Exec() error: %v", err)
		fmt.Println(errMsg) // or handle the error message in some other way
	}

	// Prepare router
	router := gin.New()

	// Register routes
	routes.RegisterRouter(ctx, pool, router)

	// Prepare server
	serveAddr := ":" + fmt.Sprint(cfg.AppPort)
	server := &http.Server{
		Addr:    serveAddr,
		Handler: router,
	}

	// Start http server
	fmt.Printf("Serving on http://localhost:%s\n", fmt.Sprint(cfg.AppPort))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("HTTP server error: %s\n", err)
	}
}
