package model

import "time"

type Participant struct {
	ID                                string                            `json:"id" dynamodbav:"ID"`
	CreatedAt                         time.Time                         `json:"createdAt"`
	UpdatedAt                         time.Time                         `json:"updatedAt"`
	DeletedAt                         *time.Time                        `json:"deletedAt"`
	DiscussionID                      string                            `json:"discussionID"`
	Discussion                        *Discussion                       `json:"discussion" dynamodbav:"-"`
	ViewerID                          string                            `json:"viewerID"`
	Viewer                            *Viewer                           `json:"viewer" dynamodbav:"-"`
	DiscussionNotificationPreferences DiscussionNotificationPreferences `json:"discussionNotificationPreferences"`
	Posts                             *PostsConnection                  `json:"posts" dynamodbav:"-"`
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
