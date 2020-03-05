package model

type DiscussionsEdge struct {
	Cursor string      `json:"cursor"`
	Node   *Discussion `json:"node"`
}
