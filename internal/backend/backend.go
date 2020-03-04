package backend

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/nedrocks/delphisbe/internal/datastore"
)

type DAOManager interface {
	GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error)
}

type daoManager struct {
	db datastore.Datastore
}

func NewDaoManager() DAOManager {
	return &daoManager{
		datastore.NewDatastore(config.TablesConfig{}),
	}
}
