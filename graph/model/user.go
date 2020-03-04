package model

import "time"

type User struct {
	ID           string                  `json:"id"`
	CreatedAt    time.Time               `json:"createdAt"`
	UpdatedAt    time.Time               `json:"updatedAt"`
	DeletedAt    *time.Time              `json:"deletedAt"`
	Participants *ParticipantsConnection `json:"participants" dynamodbav:"-"`
	Viewers      *ViewersConnection      `json:"viewers" dynamodbav:"-"`
}
