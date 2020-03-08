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
	GetParticipantByID(ctx context.Context, discussionParticipantKey model.DiscussionParticipantKey) (*model.Participant, error)
	GetParticipantsByIDs(ctx context.Context, discussionParticipantKeys []model.DiscussionParticipantKey) (map[model.DiscussionParticipantKey]*model.Participant, error)
	AddParticipantAndViewerToUser(ctx context.Context, userID string, participantID int, discussionID string, viewerID string) (*model.User, error)
	CreatePost(ctx context.Context, discussionKey model.DiscussionParticipantKey, content string) (*model.Post, error)
	GetPostsByDiscussionID(ctx context.Context, discussionID string) ([]*model.Post, error)
	AddViewerToUser(ctx context.Context, userID, discussionID, viewerID string) (*model.User, error)
	GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error)
	CreateUser(ctx context.Context) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetViewerByID(ctx context.Context, discussionViewerKey model.DiscussionViewerKey) (*model.Viewer, error)
	GetViewersByIDs(ctx context.Context, discussionViewerKeys []model.DiscussionViewerKey) (map[model.DiscussionViewerKey]*model.Viewer, error)
	CreateViewerForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Viewer, error)
}

type daoManager struct {
	db datastore.Datastore
}

func NewDaoManager(conf config.Config) DAOManager {
	return &daoManager{
		datastore.NewDatastore(conf.DBConfig),
	}
}
