package env

import (
	"errors"
	"os"
	"strconv"
	"sync"
)

var (
	envSyncOnce sync.Once

	isProduction bool
	isStaging    bool
	secretKey    string
)

// IsProduction returns true when env is set to production.
func IsProduction() bool {
	initializeEnvs()
	return isProduction
}

func IsStaging() bool {
	initializeEnvs()
	return isStaging
}

// GetSecretKey fetches the app secret key from env.
func GetSecretKey() string {
	initializeEnvs()
	return secretKey
}

func initializeEnvs() {
	envSyncOnce.Do(func() {
		isProduction = os.Getenv("ENV") == "production"
		isStaging = os.Getenv("ENV") == "staging"
		secretKey = os.Getenv("JWT_SECRET")
	})
}

func GetEnv(key string) (string, error) {
	s := os.Getenv(key)
	if s == "" {
		return s, errors.New("getenv: environment variable empty")
	}
	return s, nil
}

func GetEnvInt(key string) (int, error) {
	s, err := GetEnv(key)
	if err != nil {
		return 0, err
	}

	v, err := strconv.Atoi(s)
	return v, nil
}

func GetEnvBool(key string) (bool, error) {
	s, err := GetEnv(key)
	if err != nil {
		return false, err
	}

	v, err := strconv.ParseBool(s)
	if nil != err {
		return false, err
	}
	return v, nil
}
