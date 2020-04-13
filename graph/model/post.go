package model

import "time"

type Post struct {
	ID                string             `json:"id" dynamodbav:"ID" gorm:"type:varchar(32)"`
	CreatedAt         time.Time          `json:"createdAt" dynamodbav:"CreatedAt" gorm:"not null"`
	UpdatedAt         time.Time          `json:"updatedAt" gorm:"not null"`
	DeletedAt         *time.Time         `json:"deletedAt"`
	DeletedReasonCode *PostDeletedReason `json:"deletedReasonCode" gorm:"type:varchar(32)"`
	Discussion        *Discussion        `json:"discussion" dynamodbav:"-" gorm:"-"` //gorm:"foreignkey:DiscussionID"`
	DiscussionID      string             `json:"discussionID" dynamodbav:"DiscussionID" gorm:"type:varchar(32)"`
	Participant       *Participant       `json:"participant" dynamodbav:"-" gorm:"-"` //gorm:"foreignkey:participant_id;association_foreignkey:id"`
	ParticipantID     int                `json:"participantID"`
	PostContentID     string             `json:"postContentID" gorm:"type:varchar(32)"`
	PostContent       PostContent        `json:"postContent" gorm:"-"` //gorm:"foreinkey:post_content_id;association_foreignkey:id"`
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
