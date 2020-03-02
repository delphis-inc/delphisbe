package model

type User struct {
	ID           string                  `json:"id"`
	Participants *ParticipantsConnection `json:"participants"`
	Viewers      *ViewersConnection      `json:"viewers"`
}
