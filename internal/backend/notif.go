package backend

import (
	"context"
	"fmt"
	"sort"

	"github.com/delphis-inc/delphisbe/internal/util"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/notif"
	"github.com/sirupsen/logrus"
)

type SingleNotificationSendStatus struct {
	UserID      string
	Device      *model.UserDevice
	HasSent     bool
	HasFinished bool
	Success     bool
}

type SendNotificationResponse struct {
	ProgressChannel    <-chan *SingleNotificationSendStatus
	NumDevicesToNotify int
}

func (d *delphisBackend) SendNotificationsToSubscribers(ctx context.Context, userID string, discussion *model.Discussion, post *model.Post, contentPreview *string) (*SendNotificationResponse, error) {
	if discussion == nil || post == nil {
		return nil, fmt.Errorf("Cannot send notification to missing Post or Discussion")
	}

	if post.PostContent == nil {
		resp, err := d.db.GetPostContentByID(ctx, *post.PostContentID)
		if err != nil {
			logrus.WithError(err).Error("failed to get post content by ID")
			return nil, err
		}

		post.PostContent = resp
	}

	// Get users that want notifications on everything
	iter := d.db.GetDUAForEverythingNotifications(ctx, discussion.ID, userID)
	usersToNotify, err := d.db.DuaIterCollect(ctx, iter)
	if err != nil {
		logrus.WithError(err).Error("failed to get users to notify")
		return nil, err
	}

	// Get Users that were mentioned
	if post.PostContent.MentionedEntities != nil {
		mentionedUsers, err := d.getMentionedUsersToNotify(ctx, "", discussion.ID, post.PostContent.MentionedEntities)
		if err != nil {
			logrus.WithError(err).Error("failed to get mentioned users")
			return nil, err
		}

		usersToNotify = append(usersToNotify, mentionedUsers...)
	}

	notifChan := make(chan *SingleNotificationSendStatus, 1)
	go func() {
		// Track the users we have sent notifications to
		seenUserIDs := map[string]bool{}
		for _, user := range usersToNotify {
			if _, ok := seenUserIDs[user.UserID]; !ok {
				seenUserIDs[user.UserID] = true

				go func() {
					sendStatus := &SingleNotificationSendStatus{
						UserID:      user.UserID,
						Device:      nil,
						HasSent:     false,
						HasFinished: true,
						Success:     false,
					}
					if user == nil {
						sendMessageNonBlocking(notifChan, sendStatus)
						return
					}
					userDevices, err := d.GetUserDevicesByUserID(ctx, user.UserID)
					if err != nil {
						sendMessageNonBlocking(notifChan, sendStatus)
						return
					}
					if len(userDevices) == 0 {
						sendMessageNonBlocking(notifChan, sendStatus)
						return
					}
					sort.Slice(userDevices, func(lhs, rhs int) bool {
						return userDevices[lhs].LastSeen.Before(userDevices[rhs].LastSeen)
					})
					toSendTo := userDevices[0]
					sendStatus.Device = &toSendTo
					if toSendTo.Token == nil || len(*toSendTo.Token) == 0 {
						sendMessageNonBlocking(notifChan, sendStatus)
						return
					}
					notificationBody, err := notif.BuildPushNotification(ctx, *discussion, *post, contentPreview)
					if err != nil || notificationBody == nil {
						sendMessageNonBlocking(notifChan, sendStatus)
						return
					}
					sent, err := notif.SendPushNotification(ctx, d.config.AblyConfig, &toSendTo, *notificationBody)
					sendStatus.HasSent = sent
					sendStatus.HasFinished = true
					sendStatus.Success = err == nil
					sendMessageNonBlocking(notifChan, sendStatus)
				}()
			}
		}
	}()
	return &SendNotificationResponse{
		ProgressChannel:    notifChan,
		NumDevicesToNotify: len(discussion.Participants),
	}, nil
}

func sendMessageNonBlocking(notifChan chan *SingleNotificationSendStatus, status *SingleNotificationSendStatus) bool {
	select {
	case notifChan <- status:
		// The channel sent
		logrus.Debugf("SendNotificationsToSubscriber::Notified Channel of failure")
		return true
	default:
		// Nothing happened -- don't block.
		return false
	}
}

func (d *delphisBackend) getMentionedUsersToNotify(ctx context.Context, userID string, discussionID string, mentionedEntities []string) ([]*model.DiscussionUserAccess, error) {
	var participantIDs []string
	var userIDs []string

	// Parse mentioned entities for participants
	for _, entity := range mentionedEntities {
		parsedEntity, err := util.ReturnParsedEntityID(entity)
		if err != nil {
			logrus.WithError(err).Error("failed to parse entity")
			continue
		}
		if parsedEntity.Type == model.ParticipantPrefix {
			participantIDs = append(participantIDs, parsedEntity.ID)
		}
	}

	// Get participants by IDs so we can have the userIDa
	participants, err := d.db.GetParticipantsByIDs(ctx, participantIDs)
	if err != nil {
		logrus.WithError(err).Error("failed to get participants by IDs")
		return nil, err
	}

	// Create UserID list
	for _, participant := range participants {
		userIDs = append(userIDs, *participant.UserID)
	}

	iter := d.db.GetDUAForMentionNotifications(ctx, discussionID, userID, userIDs)
	notifyUsers, err := d.db.DuaIterCollect(ctx, iter)
	if err != nil {
		logrus.WithError(err).Error("failed to get ")
		return nil, err
	}

	return notifyUsers, nil
}
