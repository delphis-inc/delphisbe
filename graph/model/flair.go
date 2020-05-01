package model

import "time"

type Flair struct {
	ID          string         `json:"id" dynamodbav:"ID" gorm:"type:varchar(36);primary_key;"`
	TemplateID  string         `json:"templateID" dynamodbav:"TemplateID" gorm:"type:varchar(36);"`
	Template    *FlairTemplate `json:"template" dynamodbav:"-" gorm:"foreignKey:TemplateID;"`
	UserID      string         `json:"userID" dynamodbav:"UserID" gorm:"type:varchar(36);"`
	User        *User          `json:"user" dynamodbav:"-" gorm:"foreignKey:UserID;"`
	CreatedAt   time.Time      `json:"createdAt" gorm:"NOT NULL;default:CURRENT_TIMESTAMP;"`
	UpdatedAt   time.Time      `json:"updatedAt" gorm:"NOT NULL;default:CURRENT_TIMESTAMP ONUPDATE CURRENT_TIMESTAMP;"`
	DeletedAt   *time.Time     `json:"deletedAt"`
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
