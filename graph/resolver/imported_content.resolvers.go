package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/model"
)

func (r *contentQueueRecordResolver) CreatedAt(ctx context.Context, obj *model.ContentQueueRecord) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

func (r *contentQueueRecordResolver) UpdatedAt(ctx context.Context, obj *model.ContentQueueRecord) (string, error) {
	return obj.UpdatedAt.Format(time.RFC3339), nil
}

func (r *contentQueueRecordResolver) IsDeleted(ctx context.Context, obj *model.ContentQueueRecord) (bool, error) {
	return obj.DeletedAt != nil, nil
}

func (r *contentQueueRecordResolver) HasPosted(ctx context.Context, obj *model.ContentQueueRecord) (bool, error) {
	return obj.PostedAt != nil, nil
}

func (r *importedContentResolver) CreatedAt(ctx context.Context, obj *model.ImportedContent) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

func (r *importedContentResolver) MatchingTags(ctx context.Context, obj *model.ImportedContent) ([]string, error) {
	return obj.Tags, nil
}

// ContentQueueRecord returns generated.ContentQueueRecordResolver implementation.
func (r *Resolver) ContentQueueRecord() generated.ContentQueueRecordResolver {
	return &contentQueueRecordResolver{r}
}

// ImportedContent returns generated.ImportedContentResolver implementation.
func (r *Resolver) ImportedContent() generated.ImportedContentResolver {
	return &importedContentResolver{r}
}

type contentQueueRecordResolver struct{ *Resolver }
type importedContentResolver struct{ *Resolver }
