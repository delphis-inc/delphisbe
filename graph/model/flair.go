package model

type Flair struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	ImageURL    string `json:"imageURL"`
	Source      string `json:"source"`
}
