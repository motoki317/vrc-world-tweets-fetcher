package gorm

import (
	"gorm.io/gorm/clause"

	"github.com/motoki317/vrc-world-tweets-fetcher/db/model"
)

func (r *Repo) CreateWorldIfNotExists(world *model.World) (created bool, err error) {
	res := r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(world)
	return res.RowsAffected > 0, convertError(res.Error)
}
