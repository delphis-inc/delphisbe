package model

import "time"

type Participant struct {
	ParticipantID                     int                               `json:"participantID" dynamodbav:"ParticipantID"`
	CreatedAt                         time.Time                         `json:"createdAt"`
	UpdatedAt                         time.Time                         `json:"updatedAt"`
	DeletedAt                         *time.Time                        `json:"deletedAt"`
	DiscussionID                      string                            `json:"discussionID" dynamodbav:"DiscussionID"`
	Discussion                        *Discussion                       `json:"discussion" dynamodbav:"-"`
	ViewerID                          string                            `json:"viewerID"`
	Viewer                            *Viewer                           `json:"viewer" dynamodbav:"-"`
	DiscussionNotificationPreferences DiscussionNotificationPreferences `json:"discussionNotificationPreferences"`
	Posts                             *PostsConnection                  `json:"posts" dynamodbav:"-"`

	// NOTE: This is not exposed currently but keeping it here for
	// testing purposes. We will try out exposing user information one of the tests.
	UserID string `json:"userID"`
	User   *User  `json:"user" dynamodbav:"-"`
}

func (p Participant) DiscussionParticipantKey() DiscussionParticipantKey {
	return DiscussionParticipantKey{
		DiscussionID:  p.DiscussionID,
		ParticipantID: p.ParticipantID,
	}
}

type ParticipantsEdge struct {
	Cursor string       `json:"cursor"`
	Node   *Participant `json:"node"`
}

type ParticipantsConnection struct {
	ids  []string
	from int
	to   int
}

func (p *ParticipantsConnection) TotalCount() int {
	return len(p.ids)
}

func (p *ParticipantsConnection) PageInfo() PageInfo {
	from := EncodeCursor(p.from)
	to := EncodeCursor(p.to)
	return PageInfo{
		StartCursor: &from,
		EndCursor:   &to,
		HasNextPage: p.to < len(p.ids),
	}
}
