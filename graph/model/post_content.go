package model

import "time"

// This is really a placeholder rn. We will want the following:
// * Edit history and active version
// * Markup (e.g. when people are tagged having that be an ID or a token)
// * URL Wrapping
// * Ability to contain different types of posts (e.g. images, twitter cards)
type PostContent struct {
	ID                string    `json:"id" gorm:"type:varchar(36);"`
	Content           string    `json:"content" gorm:"type:text;"`
	MentionedEntities []string  `json:"mentionedEntities" gorm:"type:varchar(50)[];"`
	CreatedAt         time.Time `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP;"`
	UpdatedAt         time.Time `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP;"`
}
