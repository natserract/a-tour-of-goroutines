package config

import (
	"goroutines/pkg/env"
)

type Configuration struct {
	AppPort    int
	AppHost    string
	DBHost     string
	DBUsername string
	DBPass     string
	DBName     string
	DBPort     int
	DBParams   string
}

func GetConfig() *Configuration {
	dbHost, err := env.GetEnv("DB_HOST")
	if err != nil {
		return nil
	}

	dbUsername, err := env.GetEnv("DB_USERNAME")
	if err != nil {
		return nil
	}

	dbPass, err := env.GetEnv("DB_PASSWORD")
	if err != nil {
		return nil
	}

	dbName, err := env.GetEnv("DB_NAME")
	if err != nil {
		return nil
	}

	dbPort, err := env.GetEnvInt("DB_PORT")
	if err != nil {
		return nil
	}

	dbParams, err := env.GetEnv("DB_PARAMS")
	if err != nil {
		return nil
	}

	return &Configuration{
		AppPort:    8080,
		AppHost:    "localhost",
		DBHost:     dbHost,
		DBUsername: dbUsername,
		DBPass:     dbPass,
		DBName:     dbName,
		DBPort:     dbPort,
		DBParams:   dbParams,
	}
}
