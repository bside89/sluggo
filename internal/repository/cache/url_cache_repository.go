// Package cache provides caching mechanisms for the URL shortening service.
package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sluggo/internal/domain"

	"github.com/redis/go-redis/v9"
)

const keyPrefix = "url:"

type urlCacheRepository struct {
	client *redis.Client
}

// New returns a domain.URLCache backed by Redis.
func New(client *redis.Client) domain.URLCache {
	return &urlCacheRepository{client: client}
}

// Get retrieves the long URL associated with the given hash from the cache.
func (r *urlCacheRepository) Get(ctx context.Context, hash string) (string, error) {
	val, err := r.client.Get(ctx, keyPrefix+hash).Result()
	if errors.Is(err, redis.Nil) {
		return "", domain.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("cache get: %w", err)
	}
	return val, nil
}

// Set stores the long URL in the cache with the given hash and TTL.
func (r *urlCacheRepository) Set(ctx context.Context, hash string, longURL string, ttl time.Duration) error {
	if err := r.client.Set(ctx, keyPrefix+hash, longURL, ttl).Err(); err != nil {
		return fmt.Errorf("cache set: %w", err)
	}
	return nil
}
