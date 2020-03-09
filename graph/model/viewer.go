package model

import "time"

type Viewer struct {
	ID                      string                        `json:"id" dynamodbav:"ViewerID"`
	CreatedAt               time.Time                     `json:"createdAt"`
	UpdatedAt               time.Time                     `json:"updatedAt"`
	DeletedAt               *time.Time                    `json:"deletedAt"`
	NotificationPreferences ViewerNotificationPreferences `json:"notificationPreferences"`
	DiscussionID            string                        `json:"discussionID" dynamodbav:"DiscussionID"`
	Discussion              *Discussion                   `json:"discussion" dynamodbav:"-"`
	LastViewed              *time.Time                    `json:"lastViewed"`
	LastViewedPostID        *string                       `json:"lastViewedPostID"`
	LastViewedPost          *Post                         `json:"lastViewedPost" dynamodbav:"-"`
	Bookmarks               *PostsConnection              `json:"bookmarks" dynamodbav:"-"`

	// NOTE: This is not exposed currently but keeping it here for
	// testing purposes. We will try out exposing user information one of the tests.
	UserID string `json:"userID"`
	User   *User  `json:"user" dynamodbav:"-"`
}

func (v Viewer) DiscussionViewerKey() DiscussionViewerKey {
	return DiscussionViewerKey{
		DiscussionID: v.DiscussionID,
		ViewerID:     v.ID,
	}
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
