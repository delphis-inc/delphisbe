package model

import "time"

type UserDevice struct {
	ID        string     `json:"id" gorm:"type:varchar(36);"`
	CreatedAt time.Time  `json:"createdAt" gorm:"not null;default:CURRENT_TIMESTAMP;"`
	DeletedAt *time.Time `json:"deletedAt"`
	Platform  string     `json:"platform" gorm:"not null;"`
	LastSeen  time.Time  `json:"lastSeen" gorm:"not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;"`

	Token *string `json:"token" gorm:"type:varchar(128);"`

	UserID *string `json:"userID" gorm:"type:varchar(36);"`
	User   *User   `json:"user" gorm:"foreignKey:UserID;"`
}
