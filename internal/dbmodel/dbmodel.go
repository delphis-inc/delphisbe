package dbmodel

import (
	"time"

	"github.com/delphis-inc/delphisbe/graph/model"
)

type Discussion struct {
	ID            string `json:"id"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
	Title         string
	AnonymityType model.AnonymityType
	//Moderator Moderator
	ModeratorID string
}

type DiscussionParticipantKey struct {
	Discussion   Discussion `gorm:"foreignkey:ID;association_foreignkey:DiscussionID"`
	DiscussionID string
	//Participant Participant `gorm:"foreignkey:ParticipantID;association_foreignkey:ParticipantID`
	ParticipantID int
}
