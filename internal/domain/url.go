package domain

import "time"

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
