package model

import "time"

type ImportedContent struct {
	ID          string
	CreatedAt   time.Time
	ContentName string
	ContentType string
	Link        string
	Overview    string
	Source      string
	Tags        []string
}

type Tag struct {
	ID        string
	Tag       string
	CreatedAt time.Time
	DeletedAt *time.Time
}

type ImportedContentInput struct {
	ContentName string `json:"content_name"`
	ContentType string `json:"content_type"`
	Link        string `json:"link"`
	Overview    string `json:"overview"`
	Source      string `json:"source"`
	Tags        string `json:"tags"`
}

type ContentQueueRecord struct {
	DiscussionID      string
	ImportedContentID string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
	PostedAt          *time.Time
	MatchingTags      []string
}
