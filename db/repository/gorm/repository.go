package gorm

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/motoki317/vrc-world-tweets-fetcher/db/repository"
	"github.com/motoki317/vrc-world-tweets-fetcher/utils"
)

type Repo struct {
	db *gorm.DB
}

// Interface type guard
var _ repository.Repository = (*Repo)(nil)

func NewGormRepository() (*Repo, error) {
	host := utils.MustGetEnv("MYSQL_HOST")
	port := utils.MustGetEnv("MYSQL_PORT")
	user := utils.MustGetEnv("MYSQL_USERNAME")
	pass := utils.MustGetEnv("MYSQL_PASSWORD")
	database := utils.MustGetEnv("MYSQL_DATABASE")
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=true&loc=Local",
		user, pass, host, port, database,
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                     dsn,
		DefaultStringSize:       256,
		DontSupportRenameIndex:  true,
		DontSupportRenameColumn: true,
	}))
	if err != nil {
		return nil, err
	}
	return &Repo{db: db}, nil
}

// DB returns the underlying database connection.
func (r *Repo) DB() *gorm.DB {
	return r.db
}
