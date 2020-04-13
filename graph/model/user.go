package model

import "time"

type User struct {
	ID            string       `json:"id" dynamodbav:"ID" gorm:"type:varchar(32)"`
	CreatedAt     time.Time    `json:"createdAt" gorm:"not null"`
	UpdatedAt     time.Time    `json:"updatedAt" gorm:"not null"`
	DeletedAt     *time.Time   `json:"deletedAt"`
	UserProfileID string       `json:"userProfileID" gorm:"type:varchar(32)"`
	UserProfile   *UserProfile `json:"userProfile" dynamodbav:"-"`

	// Going through a `through` table so we can encrypt this in the future.
	Participants []*Participant `json:"participants" dynamodbav:"-" gorm:"-"` //gorm:"many2many:user_participants;"`
	Viewers      []*Viewer      `json:"viewers" dynamodbav:"-" gorm:"-"`      //gorm:"many2many:user_viewers;"`
}

func (u *User) GetParticipantKeyForDiscussionID(id string) *DiscussionParticipantKey {
	if len(u.Participants) == 0 {
		return nil
	}
	for _, participant := range u.Participants {
		if participant.DiscussionID == id {
			return &DiscussionParticipantKey{
				DiscussionID:  participant.DiscussionID,
				ParticipantID: participant.ParticipantID,
			}
		}
	}
	return nil
}
