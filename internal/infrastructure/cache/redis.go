// Package cache provides caching mechanisms for the URL shortening service.
package cache

import (
	"context"
	"log/slog"
	"time"

	"sluggo/config"

	"github.com/redis/go-redis/v9"
)

// Connect creates a Redis client, validates connectivity with a Ping and returns it.
// If Redis is unreachable, the error is logged and the client is still returned so the
// application can start and serve requests using the database as fallback.
func Connect(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr(),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		slog.Warn("warning: could not connect to Redis — cache disabled, falling back to database",
			slog.String("addr", cfg.RedisAddr()), slog.Any("error", err))
	}

	return client
}
