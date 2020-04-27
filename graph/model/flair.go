package model

import "time"

type Flair struct {
	ID          string     `json:"id" dynamodbav:"ID" gorm:"type:varchar(36);primary_key"`
	DisplayName *string    `json:"displayName" gorm:"type:varchar(128)"`
	ImageURL    *string    `json:"imageURL" gorm:"type:text"`
	Source      string     `json:"source" gorm:"type:varchar(128);NOT NULL"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"NOT NULL;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"NOT NULL;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

type FlairsEdge struct {
	Cursor string `json:"cursor"`
	Node   *Flair `json:"node"`
}

type FlairsConnection struct {
	ids   []string
	from int
	to   int
}

func (p *FlairsConnection) TotalCount() int {
	return len(p.ids)
}

func (p *FlairsConnection) PageInfo() PageInfo {
	from := EncodeCursor(p.from)
	to := EncodeCursor(p.to)
	return PageInfo{
		StartCursor: &from,
		EndCursor:   &to,
		HasNextPage: p.to < len(p.ids),
	}
}

// type UserFlair struct {
// 	UserID  string `json:"userID" gorm:"type:varchar(36);primary_key"`
// 	FlairID string `json:"flairID" gorm:"type:varchar(36);primary_key"`
// }
