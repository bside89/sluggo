package cache

import (
	"context"
	"log"
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
		log.Printf("warning: could not connect to Redis at %s: %v — cache disabled, falling back to database", cfg.RedisAddr(), err)
	}

	return client
}
