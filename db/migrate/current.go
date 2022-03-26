package migrate

import (
	"github.com/go-gormigrate/gormigrate/v2"

	"github.com/motoki317/vrc-world-tweets-fetcher/db/model"
)

// Migrations returns the list of all migrations.
func Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{}
}

// AllTables returns the list of all LATEST table models.
func AllTables() []interface{} {
	return []interface{}{
		&model.World{},
	}
}
