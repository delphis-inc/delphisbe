package model

import "time"

type Moderator struct {
	ID            string       `json:"id" gorm:"type:varchar(36)"`
	CreatedAt     time.Time    `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time    `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP"`
	DeletedAt     *time.Time   `json:"deletedAt"`
	UserProfileID *string      `json:"userProfileID" gorm:"type:varchar(36)"`
	UserProfile   *UserProfile `json:"userProfile" dynamodbav:"-" gorm:"foreignKey:UserProfileID"`
	Discussion    *Discussion `gorm:"-" dynamodbav:"-"`
}
