package model

import "time"

type Participant struct {
	ID            string      `json:"id" gorm:"type:varchar(32)"`
	ParticipantID int         `json:"participantID" dynamodbav:"ParticipantID"`
	CreatedAt     time.Time   `json:"createdAt" gorm:"not null"`
	UpdatedAt     time.Time   `json:"updatedAt" gorm:"not null"`
	DeletedAt     *time.Time  `json:"deletedAt"`
	DiscussionID  string      `json:"discussionID" dynamodbav:"DiscussionID" gorm:"type:varchar(32)"`
	Discussion    *Discussion `json:"discussion" dynamodbav:"-" gorm:"-"` //gorm:"foreignkey:discussion_id;association_foreignkey:id"`
	ViewerID      string      `json:"viewerID" gorm:"type:varchar(32)"`
	Viewer        *Viewer     `json:"viewer" dynamodbav:"-" gorm:"-"` //gorm:"foreignkey:viewer_id;association_foreignkey:id"`
	//DiscussionNotificationPreferences DiscussionNotificationPreferences `json:"discussionNotificationPreferences"`
	Posts *PostsConnection `json:"posts" dynamodbav:"-" gorm:"-"`

	// NOTE: This is not exposed currently but keeping it here for
	// testing purposes. We will try out exposing user information one of the tests.
	UserID string `json:"userID" gorm:"type:varchar(32)"`
	User   *User  `json:"user" dynamodbav:"-" gorm:"-"` //gorm:"foreignkey:user_id;assciation_foreignkey:id"`
}

func (p Participant) DiscussionParticipantKey() DiscussionParticipantKey {
	return DiscussionParticipantKey{
		DiscussionID:  p.DiscussionID,
		ParticipantID: p.ParticipantID,
	}
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
