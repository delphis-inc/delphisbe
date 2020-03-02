package model

import (
	"encoding/base64"
	"fmt"
)

type PageInfo struct {
	StartCursor *string `json:"startCursor"`
	EndCursor   *string `json:"endCursor"`
	HasNextPage bool    `json:"hasNextPage"`
}

func EncodeCursor(i int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("cursor%d", i+1)))
}
