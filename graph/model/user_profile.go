package model

import "fmt"

const (
	twitterURLFmt = "https://www.twitter.com/%s"
)

type UserProfile struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	// Handle without the `@` sign.
	TwitterHandle           string       `json:"twitterHandle"`
	ModeratedDiscussionsIDs []string     `json:"moderatedDiscussionIDs" dynamodbav:",stringset"`
	ModeratedDiscussions    []Discussion `json:"moderatedDiscussions" dynamodbav:"-"`
}

func (u *UserProfile) TwitterURL() URL {
	return URL{
		DisplayText: fmt.Sprintf("@%s", u.TwitterHandle),
		URL:         fmt.Sprintf(twitterURLFmt, u.TwitterHandle),
	}
}
