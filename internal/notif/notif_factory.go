package notif

import (
	"context"
	"fmt"

	"github.com/delphis-inc/delphisbe/graph/model"
)

// This is stupid for now but we will make it smarter
func truncateNotificationText(text string, maxLength int) string {
	runeValue := []rune(text)
	if len(runeValue) > maxLength {
		runeValue = runeValue[:maxLength]
	}
	return string(runeValue)
}

func BuildPushNotification(ctx context.Context, discussion model.Discussion, post model.Post, contentPreview *string) (*PushNotificationBody, error) {
	content := post.PostContent.Content
	if contentPreview != nil && content != *contentPreview {
		content = *contentPreview
	}
	title := truncateNotificationText(fmt.Sprintf("New post in %s", discussion.Title), 65)
	body := truncateNotificationText(content, 156)

	return &PushNotificationBody{
		Title: title,
		Body:  body,
	}, nil
}
