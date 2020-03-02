package model

type Discussion struct {
	ID            string           `json:"id"`
	AnonymityType AnonymityType    `json:"anonymityType"`
	Posts         *PostsConnection `json:"posts"`
}

type DiscussionsEdge struct {
	Cursor string      `json:"cursor"`
	Node   *Discussion `json:"node"`
}

type DiscussionsConnection struct {
	ids  []string
	from int
	to   int
}

func (d *DiscussionsConnection) TotalCount() int {
	return len(d.ids)
}

func (d *DiscussionsConnection) PageInfo() PageInfo {
	from := EncodeCursor(d.from)
	to := EncodeCursor(d.to)
	return PageInfo{
		StartCursor: &from,
		EndCursor:   &to,
		HasNextPage: d.to < len(d.ids),
	}
}
