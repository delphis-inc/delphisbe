package model

import "time"

type User struct {
	ID        string     `json:"id" dynamodbav:"ID"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`

	// NOTE: There is probably a better way to do this, but sticking
	// to this for now.
	DiscussionParticipants DiscussionParticipantKeys `json:"participantIDs" dynamodbav:"DiscussionParticipants"`
	DiscussionViewers      DiscussionViewerKeys      `json:"viewerIDs" dynamodbav:"DiscussionViewers"`

	Participants []*Participant `json:"participants" dynamodbav:"-"`
	Viewers      []*Viewer      `json:"viewers" dynamodbav:"-"`
}
