package model

import "time"

type Post struct {
	ID string `json:"id" dynamodbav:"ID"`
	// Sort Key
	CreatedAt         time.Time         `json:"createdAt" dynamodbav:"CreatedAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
	DeletedAt         *time.Time        `json:"deletedAt"`
	DeletedReasonCode PostDeletedReason `json:"deletedReason"`
	Discussion        *Discussion       `json:"discussion" dynamodbav:"-"`
	// Primary key
	DiscussionID  string       `json:"discussionID" dynamodbav:"DiscussionID"`
	Participant   *Participant `json:"participant" dynamodbav:"-"`
	ParticipantID int          `json:"participantID"`
	PostContentID string       `json:"postContentID"`
	PostContent   PostContent  `json:"postContent"`
}

type PostsEdge struct {
	Cursor string `json:"cursor"`
	Node   *Post  `json:"node"`
}

type PostsConnection struct {
	from string
	to   string
}

func (p *PostsConnection) TotalCount() int {
	return 0 //len(p.ids)
}

func (p *PostsConnection) PageInfo() PageInfo {
	// 	from := EncodeCursor(p.from)
	// 	to := EncodeCursor(p.to)
	// 	return PageInfo{
	// 		StartCursor: &from,
	// 		EndCursor:   &to,
	// 		HasNextPage: p.to < len(p.ids),
	// 	}
	return PageInfo{}
}
