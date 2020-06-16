package model

import "time"

const (
	ParticipantPrefix = "participant"
	DiscussionPrefix  = "discussion"
)

type Post struct {
	ID                string             `json:"id" dynamodbav:"ID" gorm:"type:varchar(36);"`
	CreatedAt         time.Time          `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP;"`
	UpdatedAt         time.Time          `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP;"`
	DeletedAt         *time.Time         `json:"deletedAt"`
	DeletedReasonCode *PostDeletedReason `json:"deletedReasonCode" gorm:"type:varchar(36);"`
	Discussion        *Discussion        `json:"discussion" dynamodbav:"-" gorm:"foreignkey:DiscussionID;"`
	DiscussionID      *string            `json:"discussionID" dynamodbav:"DiscussionID" gorm:"type:varchar(36);"`
	Participant       *Participant       `json:"participant" dynamodbav:"-" gorm:"foreignkey:ParticipantID;"`
	ParticipantID     *string            `json:"participantID" gorm:"varchar(36);"`
	PostContentID     *string            `json:"postContentID" gorm:"type:varchar(36);"`
	PostContent       *PostContent       `json:"postContent" gorm:"foreignkey:PostContentID;"`
	// TODO: Do we want to also log the post_content ID so that quoted text doesn't change?
	QuotedPostID      *string `json:"quotedPostID" gorm:"type:varchar(36);"`
	QuotedPost        *Post
	MediaID           *string
	ImportedContentID *string
}

type PostsEdge struct {
	Cursor string `json:"cursor"`
	Node   *Post  `json:"node"`
}

type PostsConnection struct {
	Edges    []*PostsEdge `json:"edges"`
	PageInfo PageInfo     `json:"pageInfo"`
}
