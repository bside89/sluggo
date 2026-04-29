package usecase

import (
	"fmt"
	"net/url"

	"sluggo/internal/domain"
)

// Encoder is the contract for generating unique short hash codes.
type Encoder interface {
	Encode() (string, error)
}

// URLUseCase contains the application business logic for URL shortening.
type URLUseCase struct {
	repo    domain.URLRepository
	encoder Encoder
	baseURL string
}

// New creates a URLUseCase wired with the given repository, encoder and base URL.
func New(repo domain.URLRepository, encoder Encoder, baseURL string) *URLUseCase {
	return &URLUseCase{repo: repo, encoder: encoder, baseURL: baseURL}
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
		return "", fmt.Errorf("persisting url: %w", err)
	}

	return fmt.Sprintf("%s/%s", uc.baseURL, hash), nil
}

// ResolveURL returns the long URL associated with the given hash.
// Returns domain.ErrNotFound when the hash does not exist.
func (uc *URLUseCase) ResolveURL(hash string) (string, error) {
	entity, err := uc.repo.FindByHash(hash)
	if err != nil {
		return "", err
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
