package model

import "time"

type Viewer struct {
	ID                      string                            `json:"id"`
	NotificationPreferences DiscussionNotificationPreferences `json:"notificationPreferences"`
	Discussion              *Discussion                       `json:"discussion"`
	LastViewed              *time.Time                        `json:"lastViewed"`
	LastPostViewed          *Post                             `json:"lastPostViewed"`
	Bookmarks               *PostsConnection                  `json:"bookmarks"`
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
