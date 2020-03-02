package model

import "time"

type PostBookmark struct {
	ID         string      `json:"id"`
	Discussion *Discussion `json:"discussion"`
	Post       *Post       `json:"post"`
	CreatedAt  time.Time   `json:"createdAt"`
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
