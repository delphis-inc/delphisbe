package model

import "time"

const (
	InviteTypeInvite                InviteType = "invite"
	InviteTypeAccessRequestAccepted InviteType = "access_granted"

	InviteLinkHostname string = "https://m.chatham.ai/d"
)

type InviteType string

type Discussion struct {
	Entity
	ID              string           `json:"id" dynamodbav:"ID" gorm:"type:varchar(36);"`
	CreatedAt       time.Time        `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP;"`
	UpdatedAt       time.Time        `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP;"`
	DeletedAt       *time.Time       `json:"deletedAt"`
	Title           string           `json:"title" gorm:"not null;"`
	AnonymityType   AnonymityType    `json:"anonymityType" gorm:"type:varchar(36);not null;"`
	ModeratorID     *string          `json:"moderatorID" gorm:"type:varchar(36);"`
	Moderator       *Moderator       `json:"moderator" gorm:"foreignKey:ModeratorID;"`
	Posts           []*Post          `gorm:"foreignKey:DiscussionID;"`
	PostConnections *PostsConnection `json:"posts" dynamodbav:"-" gorm:"-"`
	Participants    []*Participant   `json:"participants" dynamodbav:"-" gorm:"foreignKey:DiscussionID;"`
	AutoPost        bool             `json:"auto_post"`
	IdleMinutes     int              `json:"idle_minutes"`
	PublicAccess    bool             `json:"publicAccess"`
	IconURL         *string          `json:"icon_url"`
}

func (Discussion) IsEntity() {}

type DiscussionAutoPost struct {
	ID          string
	IdleMinutes int
}

type DiscussionFlairTemplateAccess struct {
	DiscussionID    string
	FlairTemplateID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

type DiscussionUserAccess struct {
	DiscussionID string
	UserID       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type DiscussionAccessRequest struct {
	ID           string `json:"id"`
	UserID       string
	DiscussionID string
	CreatedAt    string              `json:"createdAt"`
	UpdatedAt    string              `json:"updatedAt"`
	IsDeleted    bool                `json:"isDeleted"`
	Status       InviteRequestStatus `json:"status"`
}

type DiscussionInvite struct {
	ID                    string `json:"id"`
	UserID                string
	DiscussionID          string
	InvitingParticipantID string
	CreatedAt             string              `json:"createdAt"`
	UpdatedAt             string              `json:"updatedAt"`
	IsDeleted             bool                `json:"isDeleted"`
	Status                InviteRequestStatus `json:"status"`
	InviteType            InviteType
}

type DiscussionLinkAccess struct {
	DiscussionID      string `json:"discussionID"`
	InviteLinkSlug    string `json:"inviteLinkSlug"`
	VipInviteLinkSlug string `json:"vipInviteLinkSlug"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
	IsDeleted         bool   `json:"isDeleted"`
}
