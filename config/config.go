package config

import (
	"goroutines/pkg/env"
	"os"
)

// Container contains environment variables for the application, database, cache, token, and http server
type (
	Container struct {
		App *App
		DB  *DB
	}
	// App contains all the environment variables for the application
	App struct {
		Port int
		Host string
	}
	// Database contains all the environment variables for the database
	DB struct {
		Host     string
		Username string
		Pass     string
		Name     string
		Port     int
		Params   string
	}
)

func New() (*Container, error) {
	app := &App{
		Port: 8080,
		Host: "localhost",
	}

	port, err := env.GetEnvInt("DB_PORT")
	if err != nil {
		return nil, err
	}

	db := &DB{
		Host:     os.Getenv("DB_HOST"),
		Username: os.Getenv("DB_USERNAME"),
		Pass:     os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Port:     port,
		Params:   os.Getenv("DB_PARAMS"),
	}

	return &Container{
		app,
		db,
	}, nil
}
