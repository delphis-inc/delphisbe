package model

type Moderator struct {
	ID            string       `json:"id"`
	DiscussionID  string       `json:"discussionID"`
	Discussion    *Discussion  `json:"discussion" dynamodbav:"-"`
	UserProfileID string       `json:"userProfileID"`
	UserProfile   *UserProfile `json:"userProfile" dynamodbav:"-"`
}
