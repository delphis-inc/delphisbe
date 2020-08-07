package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/delphis-inc/delphisbe/graph/generated"
	"github.com/delphis-inc/delphisbe/graph/model"
)

func (r *viewerResolver) NotificationPreferences(ctx context.Context, obj *model.Viewer) (model.DiscussionNotificationPreferences, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *viewerResolver) Discussion(ctx context.Context, obj *model.Viewer) (*model.Discussion, error) {
	if obj.Discussion == nil && obj.DiscussionID != nil {
		res, err := r.DAOManager.GetDiscussionByID(ctx, *obj.DiscussionID)

		if err != nil {
			return nil, err
		}
		obj.Discussion = res
	}
	return obj.Discussion, nil
}

func (r *viewerResolver) LastViewedPost(ctx context.Context, obj *model.Viewer) (*model.Post, error) {
	if obj.LastViewedPostID == nil || obj.DiscussionID == nil {
		return nil, nil
	}
	if obj.LastViewedPost == nil {
		post, err := r.DAOManager.GetPostByDiscussionPostID(ctx, *obj.DiscussionID, *obj.LastViewedPostID)
		if err != nil {
			return nil, err
		}

		obj.LastViewedPost = post
	}
	return obj.LastViewedPost, nil
}

func (r *viewerResolver) Bookmarks(ctx context.Context, obj *model.Viewer) ([]*model.PostBookmark, error) {
	panic(fmt.Errorf("not implemented"))
}

// Viewer returns generated.ViewerResolver implementation.
func (r *Resolver) Viewer() generated.ViewerResolver { return &viewerResolver{r} }

type viewerResolver struct{ *Resolver }
