package model

import "time"

type Viewer struct {
	ID                      string                            `json:"id" dynamodbav:"ID"`
	CreatedAt               time.Time                         `json:"createdAt"`
	UpdatedAt               time.Time                         `json:"updatedAt"`
	DeletedAt               *time.Time                        `json:"deletedAt"`
	NotificationPreferences DiscussionNotificationPreferences `json:"notificationPreferences"`
	DiscussionID            string                            `json:"discussionID"`
	Discussion              *Discussion                       `json:"discussion"`
	LastViewed              *time.Time                        `json:"lastViewed"`
	LastViewedPostID        *string                           `json:"lastViewedPostID"`
	LastViewedPost          *Post                             `json:"lastViewedPost" dynamodbav:"-"`
	Bookmarks               *PostsConnection                  `json:"bookmarks" dynamodbav:"-"`
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
