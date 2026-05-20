// Package config provides a centralized way to load and access application settings
// from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application settings loaded from environment variables.
type Config struct {
	AppPort       string
	AppBaseURL    string
	AppEnv        string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	HashSecretKey string
	SnowflakeNode int64
	RedisHost     string
	RedisPort     string
	RedisPassword string
	CacheTTL      time.Duration
}

// Load reads the .env file (if present) and builds a Config from env vars.
func Load() (*Config, error) {
	// Non-fatal: .env may not exist in production container environments.
	_ = godotenv.Load()

	node, err := strconv.ParseInt(getEnv("SNOWFLAKE_NODE", "1"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid SNOWFLAKE_NODE value: %w", err)
	}

	cacheTTL, err := time.ParseDuration(getEnv("CACHE_TTL", "24h"))
	if err != nil {
		return nil, fmt.Errorf("invalid CACHE_TTL value: %w", err)
	}

	return &Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		AppBaseURL:    getEnv("APP_BASE_URL", "http://localhost:8080"),
		AppEnv:        getEnv("APP_ENV", "local"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "sluggo"),
		DBPassword:    getEnv("DB_PASSWORD", "sluggo_pass"),
		DBName:        getEnv("DB_NAME", "sluggo"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		HashSecretKey: getEnv("HASH_SECRET_KEY", ""),
		SnowflakeNode: node,
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		CacheTTL:      cacheTTL,
	}, nil
}

// RedisAddr returns the Redis address in host:port format.
func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
