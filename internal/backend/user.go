package backend

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
)

// AddParticipantToUser updates a user object to refer to the participant
// and optionally the viewer.
func (d *daoManager) AddParticipantToUser(ctx context.Context, userID, participantID string, includeViewer bool) error {
	err := d.db.AddParticipantToUser(ctx, userID, participantID)

	if err != nil {
		return err
	}

	if includeViewer {
		err = d.db.AddViewerToUser(ctx, userID, participantID)

		if err != nil {
			return err
		}
	}

	return nil
}

func (d *daoManager) CreateUser(ctx context.Context) (*model.User, error) {
	userObj := model.User{
		ID:        util.UUIDv4(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	resp, err := d.db.PutUser(ctx, userObj)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
