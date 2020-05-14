package notif

import (
	"context"
	"fmt"

	"github.com/nedrocks/delphisbe/graph/model"
)

// This is stupid for now but we will make it smarter
func truncateNotificationText(text string, maxLength int) string {
	runeValue := []rune(text)
	if len(runeValue) > maxLength {
		runeValue = runeValue[:maxLength]
	}
	return string(runeValue)
}

func BuildPushNotification(ctx context.Context, discussion model.Discussion, post model.Post) (*PushNotificationBody, error) {
	title := truncateNotificationText(fmt.Sprintf("New post in %s", discussion.Title), 65)
	body := truncateNotificationText(post.PostContent.Content, 156)

	return &PushNotificationBody{
		Title: title,
		Body:  body,
	}, nil
}