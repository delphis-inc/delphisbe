package model

import "time"

type Discussion struct {
	ID            string           `json:"id" dynamodbav:"ID"`
	CreatedAt     time.Time        `json:"createdAt"`
	UpdatedAt     time.Time        `json:"updatedAt"`
	DeletedAt     *time.Time       `json:"deletedAt"`
	AnonymityType AnonymityType    `json:"anonymityType"`
	Moderator     Moderator        `json:"moderator" dynamodbav:"Moderator"`
	Posts         *PostsConnection `json:"posts" dynamodbav:"-"`
	Participants  []*Participant   `json:"participants" dynamodbav:"-"`
}
