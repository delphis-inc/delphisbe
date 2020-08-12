package model

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

type DiscussionArchive struct {
	DiscussionID string         `json:"discussionID"`
	Archive      postgres.Jsonb `json:"archive"`
	CreatedAt    time.Time      `json:"createdAt"`
}
