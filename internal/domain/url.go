package domain

import (
	"context"
	"time"
)

// URL is the core domain entity representing a shortened URL record.
type URL struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Hash      string    `gorm:"uniqueIndex;size:20;not null"`
	LongURL   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// URLRepository defines the persistence contract for URL entities.
type URLRepository interface {
	Save(url *URL) error
	FindByHash(hash string) (*URL, error)
}

// URLCache defines the caching contract for resolved URLs.
type URLCache interface {
	// Get returns the long URL for the given hash.
	// Returns ErrNotFound if the entry does not exist in the cache.
	Get(ctx context.Context, hash string) (string, error)
	// Set stores the long URL for the given hash with the specified TTL.
	Set(ctx context.Context, hash string, longURL string, ttl time.Duration) error
}
