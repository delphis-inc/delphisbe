package model

type Moderator struct {
	ID            string       `json:"id" gorm:"type:varchar(32)"`
	DiscussionID  string       `json:"discussionID" gorm:"type:varchar(32)"`
	Discussion    *Discussion  `json:"discussion" dynamodbav:"-" gorm:"-"` //gorm:"foreignkey:discussion_id;association_foreignkey:moderator_id"`
	UserProfileID string       `json:"userProfileID" gorm:"type:varchar(32)"`
	UserProfile   *UserProfile `json:"userProfile" dynamodbav:"-" gorm:"-"` //gorm:"foreignkey:user_profile_id;association_foreignkey:id"`
}
