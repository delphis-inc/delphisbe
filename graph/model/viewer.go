package model

import "time"

type Viewer struct {
	ID               string      `json:"id" dynamodbav:"ViewerID" gorm:"type:varchar(36);"`
	CreatedAt        time.Time   `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP;"`
	UpdatedAt        time.Time   `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP;"`
	DeletedAt        *time.Time  `json:"deletedAt"`
	DiscussionID     *string     `json:"discussionID" dynamodbav:"DiscussionID" gorm:"type:varchar(36);"`
	Discussion       *Discussion `json:"discussion" dynamodbav:"-" gorm:"-"` //gorm:"foreignkey:discussion_id;association_foreignkey:id;"`
	LastViewed       *time.Time  `json:"lastViewed"`
	LastViewedPostID *string     `json:"lastViewedPostID" gorm:"type:varchar(36);"`
	LastViewedPost   *Post       `json:"lastViewedPost" dynamodbav:"-" gorm:"foreignKey:LastViewedPostID;"` //gorm:"foreignkey:last_post_viewed_id;association_foreignkey:id;"`

	// NOTE: This is not exposed currently but keeping it here for
	// testing purposes. We will try out exposing user information one of the tests.
	UserID *string `json:"userID" gorm:"type:varchar(36);"`
	User   *User   `json:"user" dynamodbav:"-" gorm:"-"`
}

type ViewersEdge struct {
	Cursor string  `json:"cursor"`
	Node   *Viewer `json:"node"`
}

type ViewersConnection struct {
	ids  []string
	from int
	to   int
}

func (v *ViewersConnection) TotalCount() int {
	return len(v.ids)
}

func (v *ViewersConnection) PageInfo() PageInfo {
	from := EncodeCursor(v.from)
	to := EncodeCursor(v.to)
	return PageInfo{
		StartCursor: &from,
		EndCursor:   &to,
		HasNextPage: v.to < len(v.ids),
	}
}
