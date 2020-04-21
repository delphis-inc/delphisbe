package model

import (
	"fmt"
	"time"
)

const (
	twitterURLFmt = "https://www.twitter.com/%s"
)

type UserProfile struct {
	ID          string     `json:"id" gorm:"type:varchar(32)"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP"`
	DeletedAt   *time.Time `json:"deletedAt"`
	DisplayName string     `json:"displayName" gorm:"type:varchar(256)"`
	UserID      *string    `json:"userID" dynamodbav:",omitempty" gorm:"type:varchar(32)"`
	// Handle without the `@` sign.
	TwitterHandle string `json:"twitterHandle"`

	SocialInfos []SocialInfo `json:"socialInfos" gorm:"foreignKey:UserProfileID;PRELOAD:true"`
}

type SocialInfo struct {
	CreatedAt         time.Time  `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time  `json:"updatedAt" gorm:"not null;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP"`
	DeletedAt         *time.Time `json:"deletedAt"`
	AccessToken       string     `json:"accessToken"`
	AccessTokenSecret string     `json:"accessTokenSecret"`
	// NOTE: This is the social network UserID (not Chatham)
	UserID          string `json:"userID"`
	ProfileImageURL string `json:"profileImageURL"`
	ScreenName      string `json:"screenName"`
	IsVerified      bool   `json:"isVerified"`
	Network         string `json:"network" gorm:"type:varchar(16);primary_key;auto_increment:false"`
	UserProfileID   string `json:"user_profile_id" gorm:"type:varchar(32);primary_key;auto_increment:false"`
}

func (u *UserProfile) TwitterURL() URL {
	return URL{
		DisplayText: fmt.Sprintf("@%s", u.TwitterHandle),
		URL:         fmt.Sprintf(twitterURLFmt, u.TwitterHandle),
	}
}
