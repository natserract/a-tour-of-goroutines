package main

import (
	"context"
	"fmt"
	"goroutines/config"
	"goroutines/internal/category"
	cr "goroutines/internal/category/repository"
	"goroutines/internal/product"
	pr "goroutines/internal/product/repository"
	"goroutines/pkg/database"
	"os"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/joho/godotenv"
)

/**
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
	routes.RegisterRouter(ctx, db, router)

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
*/

func main() {
	fmt.Println("------------------- TestWithCancellation -------------------")

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

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	// Connect to the database
	db, err := database.New(ctx, cfg.DB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	var wg sync.WaitGroup

	var categoryChan = make(chan category.Category)
	var productChan = make(chan product.Product)

	getCategory := func() (*category.Category, error) {
		repo := cr.NewCategoryRepository(db)
		categoryFound, err := repo.GetReferenceByName(ctx, "Acessories")
		if err != nil {
			return nil, err
		}

		return categoryFound, nil
	}
	getProduct := func() (*product.Product, error) {
		repo := pr.NewProductRepository(db)
		productFound, err := repo.GetReferenceById(ctx, uuid.FromStringOrNil("df1655ea-e352-4e94-b25d-7ff42ca32f7e"))
		if err != nil {
			return nil, err
		}

		return productFound, nil
	}

	// Get category
	wg.Add(1)
	go func() {
		defer wg.Done()

		categoryFound, err := getCategory()
		if err != nil {
			select {
			case <-ctx.Done():
				fmt.Println("access DB task1 error:", ctx.Err())
			}
		}
		categoryChan <- *categoryFound
	}()

	// Get product
	wg.Add(1)
	go func() {
		defer wg.Done()

		productFound, err := getProduct()
		if err != nil {

		}
		productChan <- *productFound
	}()

	go func() {
		wg.Wait()
	}()

	categoryResp := <-categoryChan
	productResp := <-productChan

	fmt.Println("ctx", ctx)
	fmt.Println("categoryResp", categoryResp)
	fmt.Println("productResp", productResp)
}
