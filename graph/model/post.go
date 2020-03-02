package model

type Post struct {
	ID string `json:"id"`
}

type PostsEdge struct {
	Cursor string `json:"cursor"`
	Node   *Post  `json:"node"`
}

type PostsConnection struct {
	ids  []string
	from int
	to   int
}

func (p *PostsConnection) TotalCount() int {
	return len(p.ids)
}

func (p *PostsConnection) PageInfo() PageInfo {
	from := EncodeCursor(p.from)
	to := EncodeCursor(p.to)
	return PageInfo{
		StartCursor: &from,
		EndCursor:   &to,
		HasNextPage: p.to < len(p.ids),
	}
}
