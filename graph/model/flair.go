package model

import "time"

type Flair struct {
	ID          string         `json:"id" dynamodbav:"ID" gorm:"type:varchar(36);primary_key;"`
	TemplateID  string         `json:"templateID" dynamodbav:"TemplateID" gorm:"type:varchar(36);"`
	Template    *FlairTemplate `json:"template" dynamodbav:"-" gorm:"foreignKey:TemplateID;"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP;"`
	UpdatedAt   time.Time      `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP;"`
	DeletedAt   *time.Time     `json:"deletedAt"`

	// NOTE: This is not exposed as of 05/01/2020
	UserID      string         `json:"userID" dynamodbav:"UserID" gorm:"type:varchar(36);"`
	User        *User          `json:"user" dynamodbav:"-" gorm:"foreignKey:UserID;"`
}
