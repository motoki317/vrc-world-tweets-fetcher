package model

import (
	"time"

	"github.com/motoki317/vrc-world-tweets-fetcher/utils"
)

type World struct {
	ID            utils.VRChatWorldID `gorm:"type:char(41);not null;primaryKey"`
	OriginalTweet string              `gorm:"type:varchar(100);not null"`
	CreatedAt     time.Time           `gorm:"precision:6"`
}

func (w World) TableName() string {
	return "worlds"
}
