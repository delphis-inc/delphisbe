package notif

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/config"
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

type ApnsPushNotificationRecipient struct {
	TransportType transportType `json:"transportType"`
	DeviceToken   string        `json:"deviceToken"`
}

type FcmPushNotificationRecipient struct {
	TransportType     transportType `json:"transportType"`
	RegistrationToken string        `json:"registrationToken"`
}

type ApplePushNotification struct {
	Recipient    ApnsPushNotificationRecipient `json:"recipient"`
	Notification PushNotificationBody          `json:"notification"`
}

type AndroidPushNotification struct {
	Recipient    FcmPushNotificationRecipient `json:"recipient"`
	Notification PushNotificationBody         `json:"notification"`
}

func SendPushNotification(ctx context.Context, config config.AblyConfig, device *model.UserDevice, pushBody PushNotificationBody) (bool, error) {
	if !config.Enabled {
		return true, nil
	}

	if device == nil || device.Token == nil {
		return false, fmt.Errorf("Device token must not be nil")
	}

	var body []byte = nil
	var err error = nil
	if strings.ToLower(device.Platform) == "android" {
		body, err = json.Marshal(AndroidPushNotification{
			Recipient: FcmPushNotificationRecipient{
				TransportType:     "fcm",
				RegistrationToken: *device.Token,
			},
			Notification: pushBody,
		})
	} else if strings.ToLower(device.Platform) == "ios" {
		body, err = json.Marshal(ApplePushNotification{
			Recipient: ApnsPushNotificationRecipient{
				TransportType: "apns",
				DeviceToken:   *device.Token,
			},
			Notification: pushBody,
		})
	} else {
		return false, errors.New("Unknown platform for device: " + device.Platform)
	}

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
