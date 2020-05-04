package model

import "time"

type User struct {
	ID           string         `json:"id" dynamodbav:"ID" gorm:"type:varchar(36);"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP;"`
	UpdatedAt    time.Time      `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP;"`
	DeletedAt    *time.Time     `json:"deletedAt"`
	UserProfile  *UserProfile   `json:"userProfile" dynamodbav:"-" gorm:"foreignkey:UserID;"`

	// Going through a `through` table so we can encrypt this in the future.
	Participants []*Participant `json:"participants" dynamodbav:"-" gorm:"foreignKey:UserID;"`
	Viewers      []*Viewer      `json:"viewers" dynamodbav:"-" gorm:"foreignKey:UserID;"`
	Flairs       []*Flair       `json:"flairs" dynamodbav:"-" gorm:"foreignKey:FlairID;"`
}
