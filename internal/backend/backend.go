package backend

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/nedrocks/delphisbe/internal/datastore"
)

type DAOManager interface {
	CreateNewDiscussion(ctx context.Context, anonymityType model.AnonymityType) (*model.Discussion, error)
	GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error)
	GetDiscussionsByIDs(ctx context.Context, ids []string) (map[string]*model.Discussion, error)
	ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error)
	CreateParticipantForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Participant, error)
	GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error)
	AddParticipantAndViewerToUser(ctx context.Context, userID string, participantID int, discussionID string, viewerID string) error
	AddViewerToUser(ctx context.Context, userID, discussionID, viewerID string) error
	GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error)
	CreateUser(ctx context.Context) (*model.User, error)
	CreateViewerForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Viewer, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
}

type daoManager struct {
	db datastore.Datastore
}

func NewDaoManager(conf config.Config) DAOManager {
	return &daoManager{
		datastore.NewDatastore(conf.DBConfig),
	}
}
