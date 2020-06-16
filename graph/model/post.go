package model

import "time"

const (
	ParticipantPrefix = "participant"
	DiscussionPrefix  = "discussion"

	// DripPostType tells us if the imported content was dripped manually by the mod,
	// automatically via the discussion tags and drip, OR scheduled by the mod.
	ManualDrip    DripPostType = "manual"
	AutoDrip      DripPostType = "auto"
	ScheduledDrip DripPostType = "scheduled"
)

type DripPostType string

type Post struct {
	ID                string             `json:"id" dynamodbav:"ID" gorm:"type:varchar(36);"`
	PostType          PostType           `json:"postType"`
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
	// from string
	// to   string
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
