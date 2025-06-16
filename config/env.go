package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

var requiredEnvVars = []string{
	"APP_ENV",
	"DB_PGSQL_HOST",
	"DB_PGSQL_PORT",
	"DB_PGSQL_DATABASE",
	"DB_PGSQL_USERNAME",
	"DB_PGSQL_PASSWORD",
	"PORT",
	"JWT_SECRET",
	"JWT_ACCESS_TOKEN_EXPIRATION",
	"JWT_REFRESH_TOKEN_EXPIRATION",
}

func CheckEnv() error {
	for _, envVar := range requiredEnvVars {
		if GetEnv(envVar, "") == "" {
			return fmt.Errorf("%s is required", envVar)
		}
	}
	return nil
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		valueInt, err := strconv.Atoi(value)

		if err != nil {
			log.Fatal(fmt.Sprintf("failed to parse %s", key), err)
		}

		return valueInt
	}

	return defaultValue
}
