package main

import (
	"context"
	"errors"
	"fmt"
	"goroutines/config"
	"goroutines/pkg/database"
	"goroutines/pkg/env"
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

	cfg, err := config.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load config: %v\n", err)
		os.Exit(1)
	}

	// Shared ctx
	ctx := context.Background()

	// Connect to the database
	db, err := database.New(ctx, cfg.DB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Printf("Successfully connected to database %v %s", cfg.DB, "\n")

	// Disable debug mode in production
	if env.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Check reachability
	if _, err = db.Exec(ctx, `SELECT 1`); err != nil {
		errMsg := fmt.Errorf("pool.Exec() error: %v", err)
		fmt.Println(errMsg) // or handle the error message in some other way
	}

	// Prepare router
	router := gin.New()

	// Register routes
	routes.RegisterRouter(ctx, db.Pool, router)

	// Prepare server
	serveAddr := ":" + fmt.Sprint(cfg.App.Port)
	server := &http.Server{
		Addr:    serveAddr,
		Handler: router,
	}

	// Start http server
	fmt.Printf("Serving on http://localhost:%s\n", fmt.Sprint(cfg.App.Port))
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("HTTP server error: %s\n", err)
	}
}
