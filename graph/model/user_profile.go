package model

import "fmt"

const (
	twitterURLFmt = "https://www.twitter.com/%s"
)

type UserProfile struct {
	ID          string `json:"id" gorm:"type:varchar(32)"`
	DisplayName string `json:"displayName" gorm:"type:varchar(256)"`
	UserID      string `json:"userID" dynamodbav:",omitempty" gorm:"type:varchar(32)"`
	User        User   `json:"user" gorm:"-"` //gorm:"foreignkey:user_id;association_foreignkey:id"`
	// Handle without the `@` sign.
	TwitterHandle string `json:"twitterHandle"`
	// ModeratedDiscussionIDs []string     `json:"moderatedDiscussionIDs" dynamodbav:",stringset,omitempty"`
	// ModeratedDiscussions   []Discussion `json:"moderatedDiscussions" dynamodbav:"-" gorm:"-"`

	// Twitter related fields
	TwitterInfo SocialInfo `json:"twitterInfo" gorm:"type:json"`
}

type SocialInfo struct {
	AccessToken       string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
	UserID            string `json:"userID"`
	ProfileImageURL   string `json:"profileImageURL"`
	ScreenName        string `json:"screenName"`
	IsVerified        bool   `json:"isVerified"`
}

func (u *UserProfile) TwitterURL() URL {
	return URL{
		DisplayText: fmt.Sprintf("@%s", u.TwitterHandle),
		URL:         fmt.Sprintf(twitterURLFmt, u.TwitterHandle),
	}
}
