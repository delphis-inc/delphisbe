package model

type Participant struct {
	ID                                string                            `json:"id"`
	Discussion                        *Discussion                       `json:"discussion"`
	Viewer                            *Viewer                           `json:"viewer"`
	DiscussionNotificationPreferences DiscussionNotificationPreferences `json:"discussionNotificationPreferences"`
	Posts                             *PostsConnection                  `json:"posts"`
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
