package util

import (
	uuid "github.com/satori/go.uuid"
)

func UUIDv4() string {
	return uuid.NewV4().String()
}
