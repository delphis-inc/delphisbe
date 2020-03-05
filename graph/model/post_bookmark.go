package model

import "time"

type PostBookmark struct {
	ID           string      `json:"id" dynamodbav:"ID"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	DeletedAt    *time.Time  `json:"deletedAt"`
	DiscussionID string      `json:"discussionID"`
	Discussion   *Discussion `json:"discussion" dynamodbav:"-"`
	PostID       string      `json:"postID"`
	Post         *Post       `json:"post" dynamodbav:"-"`
}

type PostBookmarksEdge struct {
	Cursor string        `json:"cursor"`
	Node   *PostBookmark `json:"node"`
}

type PostBookmarksConnection struct {
	ids  []string
	from int
	to   int
}

func (p *PostBookmarksConnection) TotalCount() int {
	return len(p.ids)
}

func (p *PostBookmarksConnection) PageInfo() PageInfo {
	from := EncodeCursor(p.from)
	to := EncodeCursor(p.to)
	return PageInfo{
		StartCursor: &from,
		EndCursor:   &to,
		HasNextPage: p.to < len(p.ids),
	}
}
