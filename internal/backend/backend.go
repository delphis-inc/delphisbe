package backend

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/nedrocks/delphisbe/internal/datastore"
)

type DelphisBackend interface {
	CreateNewDiscussion(ctx context.Context, creatingUser *model.User, anonymityType model.AnonymityType, title string) (*model.Discussion, error)
	GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error)
	GetDiscussionsByIDs(ctx context.Context, ids []string) (map[string]*model.Discussion, error)
	GetDiscussionByModeratorID(ctx context.Context, moderatorID string) (*model.Discussion, error)
	ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error)
	GetModeratorByID(ctx context.Context, id string) (*model.Moderator, error)
	CreateParticipantForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Participant, error)
	GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error)
	GetParticipantByID(ctx context.Context, id string) (*model.Participant, error)
	CreatePost(ctx context.Context, discussionID string, participantID string, content string) (*model.Post, error)
	GetPostsByDiscussionID(ctx context.Context, discussionID string) ([]*model.Post, error)
	GetPostContentByID(ctx context.Context, id string) (*model.PostContent, error)
	GetUserProfileByID(ctx context.Context, id string) (*model.UserProfile, error)
	GetUserProfileByUserID(ctx context.Context, userID string) (*model.UserProfile, error)
	CreateUser(ctx context.Context) (*model.User, error)
	GetOrCreateUser(ctx context.Context, input LoginWithTwitterInput) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetViewerByID(ctx context.Context, viewerID string) (*model.Viewer, error)
	GetViewersByIDs(ctx context.Context, viewerIDs []string) (map[string]*model.Viewer, error)
	CreateViewerForDiscussion(ctx context.Context, discussionID string, userID string) (*model.Viewer, error)
	GetSocialInfosByUserProfileID(ctx context.Context, userProfileID string) ([]model.SocialInfo, error)
	UpsertSocialInfo(ctx context.Context, socialInfo model.SocialInfo) (*model.SocialInfo, error)

	NewAccessToken(ctx context.Context, userID string) (*auth.DelphisAccessToken, error)
	ValidateAccessToken(ctx context.Context, token string) (*auth.DelphisAuthedUser, error)
	ValidateRefreshToken(ctx context.Context, token string) (*auth.DelphisRefreshTokenUser, error)
}

type delphisBackend struct {
	db   datastore.Datastore
	auth auth.DelphisAuth
}

func NewDelphisBackend(conf config.Config, awsSession *session.Session) DelphisBackend {
	return &delphisBackend{
		db:   datastore.NewDatastore(conf, awsSession),
		auth: auth.NewDelphisAuth(&conf.Auth),
	}
}
