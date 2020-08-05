// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
package resolver

import (
	"context"
	"fmt"
	"time"

	"github.com/delphis-inc/delphisbe/internal/cache"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/backend"
)

type resolverKeyType string

const cachedValueKey resolverKeyType = "operationCacheKey"

type Resolver struct {
	DAOManager backend.DelphisBackend
}

func GenerateCachedOperationContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, cachedValueKey, cache.NewNonDeletedInMemoryCache())
}

func GetOperationCache(ctx context.Context) cache.ChathamCache {
	cacheVal := ctx.Value(cachedValueKey)
	if cacheVal == nil {
		return nil
	}
	return cacheVal.(cache.ChathamCache)
}

func DiscussionCacheKey(id string) string {
	return fmt.Sprintf("discussion-%s", id)
}

func (r *Resolver) resolveDiscussionByID(ctx context.Context, id string) (*model.Discussion, error) {
	var discussionObj *model.Discussion

	inMemoryCache := GetOperationCache(ctx)
	if inMemoryCache != nil {
		resp, found := inMemoryCache.Get(DiscussionCacheKey(id))
		if found {
			discussionObj = resp.(*model.Discussion)
		}
	}

	if discussionObj == nil {
		var err error
		discussionObj, err = r.DAOManager.GetDiscussionByID(ctx, id)

		if err != nil {
			return nil, err
		} else if discussionObj != nil {
			inMemoryCache.Set(DiscussionCacheKey(id), discussionObj, time.Minute)
		}
	}

	return discussionObj, nil
}
