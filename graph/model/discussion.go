package model

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

const (
	InviteTypeInvite                InviteType = "invite"
	InviteTypeAccessRequestAccepted InviteType = "access_granted"

	AccessSlugLength = 12

	InviteLinkHostname string = "https://m.chatham.ai/d"
)

type InviteType string

type Discussion struct {
	Entity
	ID                    string                       `json:"id" dynamodbav:"ID" gorm:"type:varchar(36);"`
	CreatedAt             time.Time                    `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP;"`
	UpdatedAt             time.Time                    `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP;"`
	DeletedAt             *time.Time                   `json:"deletedAt"`
	Title                 string                       `json:"title" gorm:"not null;"`
	Description           string                       `json:"description"`
	TitleHistory          postgres.Jsonb               `json:"titleHistory" gorm:"type:jsonb"`
	DescriptionHistory    postgres.Jsonb               `json:"descriptionHistory" gorm:"type:jsonb"`
	AnonymityType         AnonymityType                `json:"anonymityType" gorm:"type:varchar(36);not null;"`
	ModeratorID           *string                      `json:"moderatorID" gorm:"type:varchar(36);"`
	Moderator             *Moderator                   `json:"moderator" gorm:"foreignKey:ModeratorID;"`
	Posts                 []*Post                      `gorm:"foreignKey:DiscussionID;"`
	PostConnections       *PostsConnection             `json:"posts" dynamodbav:"-" gorm:"-"`
	Participants          []*Participant               `json:"participants" dynamodbav:"-" gorm:"foreignKey:DiscussionID;"`
	AutoPost              bool                         `json:"auto_post"`
	IdleMinutes           int                          `json:"idle_minutes"`
	IconURL               *string                      `json:"icon_url"`
	DiscussionJoinability DiscussionJoinabilitySetting `json:"discussion_joinability"`
	LastPostID            *string                      `json:"last_post_id"`
	LastPostCreatedAt     *time.Time                   `json:"last_post_created_at"`
	ShuffleCount          int                          `json:"shuffle_count"`
}

type DiscussionInput struct {
	AnonymityType         *AnonymityType                `json:"anonymityType"`
	Title                 *string                       `json:"title"`
	Description           *string                       `json:"description"`
	AutoPost              *bool                         `json:"autoPost"`
	IdleMinutes           *int                          `json:"idleMinutes"`
	PublicAccess          *bool                         `json:"publicAccess"`
	IconURL               *string                       `json:"iconURL"`
	DiscussionJoinability *DiscussionJoinabilitySetting `json:"discussionJoinability"`
	LastPostID            *string                       `json:"lastPostID"`
	LastPostCreatedAt     *time.Time                    `json:"lastPostCreatedAt"`
}

type HistoricalString struct {
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
}

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
	State        DiscussionUserAccessState
	NotifSetting DiscussionUserNotificationSetting
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	RequestID    *string
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

type DiscussionAccessLink struct {
	DiscussionID string `json:"discussionID"`
	LinkSlug     string `json:"linkSlug"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func (Discussion) IsEntity() {}

func (d *Discussion) TitleHistoryAsObject() ([]*HistoricalString, error) {
	resp := make([]*HistoricalString, 0)
	history := []HistoricalString{}
	err := json.Unmarshal(d.TitleHistory.RawMessage, &history)
	if err != nil {
		// Returning an empty list for now. We may want to have error checking here
		// in the future.
		return resp, nil
	}

	for idx := range history {
		resp = append(resp, &history[idx])
	}

	return resp, nil
}

func (d *Discussion) DescriptionHistoryAsObject() ([]*HistoricalString, error) {
	resp := make([]*HistoricalString, 0)
	history := []HistoricalString{}
	err := json.Unmarshal(d.DescriptionHistory.RawMessage, &history)
	if err != nil {
		// Returning an empty list for now. We may want to have error checking here
		// in the future.
		return resp, nil
	}

	for idx := range history {
		resp = append(resp, &history[idx])
	}

	return resp, nil
}

func (d *Discussion) AddTitleToHistory(title string) error {
	historyObj, err := d.TitleHistoryAsObject()
	if err != nil {
		return err
	}

	historyObj = append(historyObj, &HistoricalString{
		Value:     title,
		CreatedAt: time.Now(),
	})

	marshaled, err := json.Marshal(historyObj)
	if err != nil {
		return err
	}

	d.TitleHistory = postgres.Jsonb{RawMessage: marshaled}

	return nil
}

func (d *Discussion) AddDescriptionToHistory(description string) error {
	historyObj, err := d.DescriptionHistoryAsObject()
	if err != nil {
		return err
	}

	historyObj = append(historyObj, &HistoricalString{
		Value:     description,
		CreatedAt: time.Now(),
	})

	marshaled, err := json.Marshal(historyObj)
	if err != nil {
		return err
	}

	d.DescriptionHistory = postgres.Jsonb{RawMessage: marshaled}

	return nil
}
