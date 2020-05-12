package backend

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *delphisBackend) UpsertUserDevice(ctx context.Context, deviceID string, userID *string, platform string, token *string) (*model.UserDevice, error) {
	userDevice := model.UserDevice{
		ID:       deviceID,
		Platform: platform,
		LastSeen: time.Now(),
		Token:    token,
		UserID:   userID,
	}

	return d.db.UpsertUserDevice(ctx, userDevice)
}

func (d *delphisBackend) GetUserDeviceByUserIDPlatform(ctx context.Context, userID string, platform string) (*model.UserDevice, error) {
	devices, err := d.db.GetUserDevicesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		// This works because we sort by LastSeen so it is the most recent one.
		if device.Platform == platform {
			return &device, nil
		}
	}
	return nil, nil
}

func (d *delphisBackend) GetUserDevicesByUserID(ctx context.Context, userID string) ([]model.UserDevice, error) {
	devices, err := d.db.GetUserDevicesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return devices, nil
}
