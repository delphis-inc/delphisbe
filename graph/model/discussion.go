package model

import "time"

type Discussion struct {
	Entity
	ID              string           `json:"id" dynamodbav:"ID" gorm:"type:varchar(36);"`
	CreatedAt       time.Time        `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP;"`
	UpdatedAt       time.Time        `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP;"`
	DeletedAt       *time.Time       `json:"deletedAt"`
	Title           string           `json:"title" gorm:"not null;"`
	AnonymityType   AnonymityType    `json:"anonymityType" gorm:"type:varchar(36);not null;"`
	ModeratorID     *string          `json:"moderatorID" gorm:"type:varchar(36);"`
	Moderator       *Moderator       `json:"moderator" gorm:"foreignKey:ModeratorID;"`
	Posts           []*Post          `gorm:"foreignKey:DiscussionID;"`
	PostConnections *PostsConnection `json:"posts" dynamodbav:"-" gorm:"-"`
	Participants    []*Participant   `json:"participants" dynamodbav:"-" gorm:"foreignKey:DiscussionID;"`
	AutoPost        bool             `json:"auto_post"`
	IdleMinutes     int              `json:"idle_minutes"`
}

func (Discussion) IsEntity() {}

type DiscussionAutoPost struct {
	ID          string
	IdleMinutes int
}
