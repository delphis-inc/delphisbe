package model

import "time"

type Discussion struct {
	ID            string           `json:"id" dynamodbav:"ID" gorm:"type:varchar(32)"`
	CreatedAt     time.Time        `json:"createdAt" gorm:"not null"`
	UpdatedAt     time.Time        `json:"updatedAt" gorm:"not null"`
	DeletedAt     *time.Time       `json:"deletedAt"`
	Title         string           `json:"title" gorm:"not null"`
	AnonymityType AnonymityType    `json:"anonymityType" gorm:"type:varchar(32);not null"`
	ModeratorID   string           `json:"moderatorID" gorm:"type:varchar(32)"`
	Moderator     Moderator        `json:"moderator" gorm:"-"` // gorm:"foreignkey:moderator_id;association_foreignkey:id"`
	Posts         *PostsConnection `json:"posts" dynamodbav:"-" gorm:"-"`
	Participants  []*Participant   `json:"participants" dynamodbav:"-" gorm:"-"`
}
