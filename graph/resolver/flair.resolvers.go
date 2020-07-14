package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/delphis-inc/delphisbe/graph/generated"
	"github.com/delphis-inc/delphisbe/graph/model"
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
			if err == nil {
				// In order to ensure we return an error craft one
				err = fmt.Errorf("No flair template found for requested flair")
			}
			return "", err
		}
		obj.Template = template
	}
	return obj.Template.Source, nil
}

// Flair returns generated.FlairResolver implementation.
func (r *Resolver) Flair() generated.FlairResolver { return &flairResolver{r} }

type flairResolver struct{ *Resolver }
