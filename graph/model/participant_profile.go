package model

type ParticipantProfile struct {
	IsAnonymous   *bool          `json:"isAnonymous"`
	Flair         *Flair         `json:"flair"`
	GradientColor *GradientColor `json:"gradientColor"`
}
