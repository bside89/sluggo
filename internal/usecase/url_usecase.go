// Package usecase provides the application business logic for the URL shortening
// service.
package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"sluggo/internal/domain"
)

// Encoder is the contract for generating unique short hash codes.
type Encoder interface {
	Encode() (string, error)
}

// URLUseCase contains the application business logic for URL shortening.
type URLUseCase struct {
	repo     domain.URLRepository
	cache    domain.URLCache
	encoder  Encoder
	baseURL  string
	cacheTTL time.Duration
}

// New creates a URLUseCase wired with the given repository, cache, encoder, base URL and cache TTL.
func New(repo domain.URLRepository, cache domain.URLCache, encoder Encoder, baseURL string, cacheTTL time.Duration) *URLUseCase {
	return &URLUseCase{
		repo:     repo,
		cache:    cache,
		encoder:  encoder,
		baseURL:  baseURL,
		cacheTTL: cacheTTL,
	}
}

// ShortenURL validates rawURL, encodes a Snowflake ID into a short hash,
// persists the mapping, and returns the full short URL.
func (uc *URLUseCase) ShortenURL(rawURL string) (string, error) {
	if err := validateURL(rawURL); err != nil {
		return "", err
	}

	hash, err := uc.encoder.Encode()
	if err != nil {
		return "", fmt.Errorf("generating short code: %w", err)
	}

	entity := &domain.URL{
		Hash:    hash,
		LongURL: rawURL,
	}
	if err := uc.repo.Save(entity); err != nil {
		slog.Error("persisting url", slog.String("hash", hash), slog.Any("error", err))
		return "", fmt.Errorf("persisting url: %w", err)
	}

	return fmt.Sprintf("%s/%s", uc.baseURL, hash), nil
}

// ResolveURL returns the long URL associated with the given hash.
// It applies the cache-aside pattern: checks the cache first, falls back to the
// database on a miss, and populates the cache with the result.
// Returns domain.ErrNotFound when the hash does not exist.
func (uc *URLUseCase) ResolveURL(hash string) (string, error) {
	ctx := context.Background()

	longURL, err := uc.cache.Get(ctx, hash)
	if err == nil {
		return longURL, nil
	}
	if !errors.Is(err, domain.ErrNotFound) {
		slog.Error("cache get error — falling back to database",
			slog.String("hash", hash), slog.Any("error", err))
	}

	entity, err := uc.repo.FindByHash(hash)
	if err != nil {
		return "", err
	}

	if err := uc.cache.Set(ctx, hash, entity.LongURL, uc.cacheTTL); err != nil {
		slog.Error("cache set error",
			slog.String("hash", hash), slog.Any("error", err))
	}

	return entity.LongURL, nil
}

func validateURL(rawURL string) error {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("invalid url: scheme must be http or https")
	}
	if u.Host == "" {
		return fmt.Errorf("invalid url: missing host")
	}
	return nil
}
