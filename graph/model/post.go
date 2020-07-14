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

	AppActionCopyToClipboard AppActionID = "db5fd0da-d645-4aa2-990c-b61d004a45e1"
	AppActionRenameChat      AppActionID = "d81118d6-427a-4267-96be-45cadd94b782"

	MutationUpdateFlairAccessToDiscussion MutationID = "4e960003-da38-4971-a23b-98953cb5ce4b"
	MutationUpdateInvitationApproval      MutationID = "84c0e197-6394-4b9a-87dc-91e75e7faf67"
	MutationUpdateViewerAccessibility     MutationID = "e8c71b3c-b984-4090-b032-7dbfd374e8c9"
	MutationUpdateDiscussionNameAndEmoji  MutationID = "633cb21f-a004-45d4-b4e8-bd6cd0bdaea9"
)

type DripPostType string
type AppActionID string
type MutationID string

type Post struct {
	DiscussionSubscriptionEntity
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
	ConciergeContent  *ConciergeContent
}

type PostsEdge struct {
	Cursor string `json:"cursor"`
	Node   *Post  `json:"node"`
}

type PostsConnection struct {
	Edges    []*PostsEdge `json:"edges"`
	PageInfo PageInfo     `json:"pageInfo"`
}
