// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
package resolver

import (
	"context"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *flairResolver) DisplayName(ctx context.Context, obj *model.Flair) (*string, error) {
	if obj.Template == nil {
		template, err := r.DAOManager.GetFlairTemplateByID(ctx, obj.TemplateID)
		if err != nil || template == nil {
			return nil, err
		}
		obj.Template = template
	}
	return obj.Template.DisplayName, nil
}

func (r *flairResolver) ImageURL(ctx context.Context, obj *model.Flair) (*string, error) {
	if obj.Template == nil {
		template, err := r.DAOManager.GetFlairTemplateByID(ctx, obj.TemplateID)
		if err != nil || template == nil {
			return nil, err
		}
		obj.Template = template
	}
	return obj.Template.ImageURL, nil
}

func (r *flairResolver) Source(ctx context.Context, obj *model.Flair) (string, error) {
	if obj.Template == nil {
		template, err := r.DAOManager.GetFlairTemplateByID(ctx, obj.TemplateID)
		if err != nil || template == nil {
			// TODO: Not sure what to return here... This really should never happen.
			return "error", err
		}
		obj.Template = template
	}
	return obj.Template.Source, nil
}

func (r *Resolver) Flair() generated.FlairResolver { return &flairResolver{r} }

type flairResolver struct{ *Resolver }
