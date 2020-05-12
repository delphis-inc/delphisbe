package notif

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/sirupsen/logrus"
)

type transportType string

const (
	ApnsTransportType transportType = "apns"
	FcmTransportType  transportType = "fcm"
)

type PushNotificationBody struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PushNotificationRecipient struct {
	TransportType transportType `json:"transportType"`
	DeviceToken   string        `json:"deviceToken"`
}

type ApplePushNotification struct {
	Recipient    PushNotificationRecipient `json:"recipient"`
	Notification PushNotificationBody      `json:"notification"`
}

func SendPushNotification(ctx context.Context, config config.AblyConfig, device *model.UserDevice, pushBody PushNotificationBody) (bool, error) {
	if !config.Enabled {
		return true, nil
	}
	body, err := json.Marshal(ApplePushNotification{
		Recipient: PushNotificationRecipient{
			TransportType: "apns",
			DeviceToken:   *device.Token,
		},
		Notification: pushBody,
	})
	if err != nil {
		return false, err
	}
	logrus.Infof("%s", string(body))
	resp, err := http.Post(fmt.Sprintf("https://%s:%s@rest.ably.io/push/publish", config.Username, config.Password), "application/json", bytes.NewBuffer(body))
	if err != nil {
		logrus.WithError(err).Errorf("Failed to post")
		return false, err
	}
	logrus.Infof("%+v", resp)
	return true, nil
}
