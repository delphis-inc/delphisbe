package model

import "time"

type FlairTemplate struct {
	ID          string     `json:"id" dynamodbav:"ID" gorm:"type:varchar(36);primary_key;"`
	DisplayName *string    `json:"displayName" gorm:"type:varchar(64);"`
	ImageURL    *string    `json:"imageURL" gorm:"type:text;"`
	Source      string     `json:"source" gorm:"type:varchar(128);not null;"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP;"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP;"`
	DeletedAt   *time.Time `json:"deletedAt"`
}
