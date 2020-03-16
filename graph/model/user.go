package model

import "time"

type User struct {
	ID            string       `json:"id" dynamodbav:"ID"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
	DeletedAt     *time.Time   `json:"deletedAt"`
	UserProfileID string       `json:"userProfileID"`
	UserProfile   *UserProfile `json:"userProfile" dynamodbav:"-"`

	// NOTE: There is probably a better way to do this, but sticking
	// to this for now.
	DiscussionParticipants *DiscussionParticipantKeys `json:"participantIDs" dynamodbav:"DiscussionParticipants,omitempty"`
	DiscussionViewers      *DiscussionViewerKeys      `json:"viewerIDs" dynamodbav:"DiscussionViewers,omitempty"`

	Participants []*Participant `json:"participants" dynamodbav:"-"`
	Viewers      []*Viewer      `json:"viewers" dynamodbav:"-"`
}

func (u *User) GetParticipantKeyForDiscussionID(id string) *DiscussionParticipantKey {
	if u.DiscussionParticipants == nil {
		return nil
	}
	for _, dpk := range u.DiscussionParticipants.Keys {
		if dpk.DiscussionID == id {
			return &dpk
		}
	}
	return nil
}
