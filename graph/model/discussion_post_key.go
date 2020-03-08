package model

import "time"

type DiscussionPostKey struct {
	DiscussionID  string    `json:"discussionID"`
	PostCreatedAt time.Time `json:"postCreatedAt"`
}
