package repository

import (
	"github.com/motoki317/vrc-world-tweets-fetcher/db/model"
)

type WorldRepository interface {
	// CreateWorldIfNotExists creates a new world record, if not exists.
	//
	// Returns created true if newly created, false otherwise.
	CreateWorldIfNotExists(world *model.World) (created bool, err error)
}
