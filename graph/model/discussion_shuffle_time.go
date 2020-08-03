package model

import "time"

type DiscussionShuffleTime struct {
	DiscussionID string     `json:"discussionID"`
	ShuffleTime  *time.Time `json:"shuffleTime"`
}
