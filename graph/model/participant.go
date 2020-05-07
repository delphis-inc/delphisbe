package model

import "time"

type Participant struct {
	ID            string           `json:"id" gorm:"type:varchar(36);"`
	ParticipantID int              `json:"participantID" dynamodbav:"ParticipantID"`
	CreatedAt     time.Time        `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP;"`
	UpdatedAt     time.Time        `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP;"`
	DeletedAt     *time.Time       `json:"deletedAt"`
	DiscussionID  *string          `json:"discussionID" dynamodbav:"DiscussionID" gorm:"type:varchar(36);"`
	Discussion    *Discussion      `json:"discussion" dynamodbav:"-" gorm:"foreignKey:DiscussionID;"`
	ViewerID      *string          `json:"viewerID" gorm:"type:varchar(36);"`
	Viewer        *Viewer          `json:"viewer" dynamodbav:"-" gorm:"foreignKey:ViewerID;"`
	Posts         *PostsConnection `json:"posts" dynamodbav:"-"`
	FlairID       *string          `json:"flairID" dynamodbav:"FlairID" gorm:"type:varchar(36);"`
	Flair         *Flair           `json:"flair" dynamodbav:"-" gorm:"foreignKey:FlairID;"`
	GradientColor *GradientColor   `json:"gradientColor" gorm:"type:varchar(36);not null;"`

	UserID *string `json:"userID" gorm:"type:varchar(36);"`
	User   *User   `json:"user" dynamodbav:"-" gorm:"foreignKey:UserID;"`

	HasJoined   bool `json:"hasJoined" gorm:"type:boolean;"`
	IsAnonymous bool `json:"isAnonymous"`
}

type ParticipantsEdge struct {
	Cursor string       `json:"cursor"`
	Node   *Participant `json:"node"`
}

type ParticipantsConnection struct {
	ids  []string
	from int
	to   int
}

func (p *ParticipantsConnection) TotalCount() int {
	return len(p.ids)
}

func (p *ParticipantsConnection) PageInfo() PageInfo {
	from := EncodeCursor(p.from)
	to := EncodeCursor(p.to)
	return PageInfo{
		StartCursor: &from,
		EndCursor:   &to,
		HasNextPage: p.to < len(p.ids),
	}
}
