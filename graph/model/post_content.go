package model

// This is really a placeholder rn. We will want the following:
// * Edit history and active version
// * Markup (e.g. when people are tagged having that be an ID or a token)
// * URL Wrapping
// * Ability to contain different types of posts (e.g. images, twitter cards)
type PostContent struct {
	ID      string `json:"id" gorm:"type:varchar(32);"`
	Content string `json:"content"`
}
