package gorm

import (
	"gorm.io/gorm"

	"github.com/motoki317/vrc-world-tweets-fetcher/db/repository"
)

// convertError converts gorm error to repository error.
func convertError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return repository.ErrNotFound
	default:
		return err
	}
}
