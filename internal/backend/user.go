package backend

import (
	"context"
	"time"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/util"
)

func (d *delphisBackend) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	return d.db.GetUserByID(ctx, userID)
}

func (d *delphisBackend) CreateUser(ctx context.Context) (*model.User, error) {
	userObj := model.User{
		ID:        util.UUIDv4(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	resp, err := d.db.UpsertUser(ctx, userObj)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
