package backend

import (
	"context"
	"fmt"
	"sort"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/notif"
	"github.com/sirupsen/logrus"
)

type SingleNotificationSendStatus struct {
	Participant model.Participant
	Device      *model.UserDevice
	HasSent     bool
	HasFinished bool
	Success     bool
}

type SendNotificationResponse struct {
	ProgressChannel    <-chan *SingleNotificationSendStatus
	NumDevicesToNotify int
}

func (d *delphisBackend) SendNotificationsToSubscribers(ctx context.Context, discussion *model.Discussion, post *model.Post, contentPreview *string) (*SendNotificationResponse, error) {
	if discussion == nil || post == nil {
		return nil, fmt.Errorf("Cannot send notification to missing Post or Discussion")
	}
	if discussion.Participants == nil {
		participants, err := d.GetParticipantsByDiscussionID(ctx, discussion.ID)
		if err != nil {
			return nil, err
		}
		discussion.Participants = make([]*model.Participant, 0)
		seenUserIDs := map[string]bool{}
		for i := range participants {
			if participants[i].UserID == nil {
				continue
			}
			userID := *participants[i].UserID
			if _, ok := seenUserIDs[userID]; !ok {
				seenUserIDs[userID] = true
				discussion.Participants = append(discussion.Participants, &participants[i])
			}
		}
	}
	notifChan := make(chan *SingleNotificationSendStatus, 1)
	go func() {
		for i := range discussion.Participants {
			participant := discussion.Participants[i]
			go func() {
				sendStatus := &SingleNotificationSendStatus{
					Participant: *participant,
					Device:      nil,
					HasSent:     false,
					HasFinished: true,
					Success:     false,
				}
				if participant.UserID == nil {
					sendMessageNonBlocking(notifChan, sendStatus)
					return
				}
				userDevices, err := d.GetUserDevicesByUserID(ctx, *participant.UserID)
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
