package postgresrepo

import (
	"errors"

	"sluggo/internal/domain"

	"gorm.io/gorm"
)

type urlRepository struct {
	db *gorm.DB
}

// New returns a domain.URLRepository backed by PostgreSQL via GORM.
func New(db *gorm.DB) domain.URLRepository {
	return &urlRepository{db: db}
}

// Save inserts a new URL record into the database.
func (r *urlRepository) Save(url *domain.URL) error {
	return r.db.Create(url).Error
}

// FindByHash retrieves a URL record by its hash.
func (r *urlRepository) FindByHash(hash string) (*domain.URL, error) {
	var url domain.URL
	err := r.db.Where("hash = ?", hash).First(&url).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &url, nil
}
