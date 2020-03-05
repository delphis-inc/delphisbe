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
	ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error)
}

type daoManager struct {
	db datastore.Datastore
}

func NewDaoManager(conf config.Config) DAOManager {
	return &daoManager{
		datastore.NewDatastore(conf.DBConfig),
	}
}
