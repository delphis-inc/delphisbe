package model

import "fmt"

const (
	twitterURLFmt = "https://www.twitter.com/%s"
)

type UserProfile struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	UserID      string `json:"userID" dynamodbav:",omitempty"`
	// Handle without the `@` sign.
	TwitterHandle          string       `json:"twitterHandle"`
	ModeratedDiscussionIDs []string     `json:"moderatedDiscussionIDs" dynamodbav:",stringset,omitempty"`
	ModeratedDiscussions   []Discussion `json:"moderatedDiscussions" dynamodbav:"-"`

	// Twitter related fields
	TwitterInfo SocialInfo `json:"twitterInfo"`
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
