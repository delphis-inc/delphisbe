// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"
	"fmt"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *participantResolver) Discussion(ctx context.Context, obj *model.Participant) (*model.Discussion, error) {
	if obj.Discussion == nil && obj.DiscussionID != nil {
		disc, err := r.resolveDiscussionByID(ctx, *obj.DiscussionID)
		if err != nil {
			return nil, err
		}
		obj.Discussion = disc
	}

	return obj.Discussion, nil
}

func (r *participantResolver) Viewer(ctx context.Context, obj *model.Participant) (*model.Viewer, error) {
	if obj.Viewer == nil && obj.ViewerID != nil {
		viewerObj, err := r.DAOManager.GetViewerByID(ctx, *obj.ViewerID)

		if err != nil {
			return nil, err
		}

		obj.Viewer = viewerObj
	}
	return obj.Viewer, nil
}

func (r *participantResolver) DiscussionNotificationPreferences(ctx context.Context, obj *model.Participant) (model.DiscussionNotificationPreferences, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *participantResolver) Posts(ctx context.Context, obj *model.Participant) ([]*model.Post, error) {
	if obj.DiscussionID == nil {
		return nil, nil
	}
	posts, err := r.DAOManager.GetPostsByDiscussionID(ctx, *obj.DiscussionID)
	if err != nil {
		return nil, err
	}

	n := 0
	for _, post := range posts {
		if post.ParticipantID != nil {
			if *post.ParticipantID == obj.ID {
				posts[n] = post
				n++
			}
		}
	}
	posts = posts[:n]
	return posts, nil
}

func (r *Resolver) Participant() generated.ParticipantResolver { return &participantResolver{r} }

type participantResolver struct{ *Resolver }
