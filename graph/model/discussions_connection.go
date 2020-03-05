package model

type DiscussionsConnection struct {
	IDs   []string
	From  int
	To    int
	Edges []*DiscussionsEdge
}

func (d *DiscussionsConnection) TotalCount() int {
	return len(d.IDs)
}

func (d *DiscussionsConnection) PageInfo() PageInfo {
	from := EncodeCursor(d.From)
	to := EncodeCursor(d.To)
	return PageInfo{
		StartCursor: &from,
		EndCursor:   &to,
		HasNextPage: d.To < len(d.IDs),
	}
}
