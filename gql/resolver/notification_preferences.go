package resolver

import (
	"github.com/graph-gophers/graphql-go"
)

type notificationPreferencesResolver struct {
}

type discussionNotificationPreferencesResolver struct {
	discussionId graphql.ID
}
