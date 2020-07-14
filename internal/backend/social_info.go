package backend

import (
	"context"

	"github.com/delphis-inc/delphisbe/graph/model"
)

func (b *delphisBackend) UpsertSocialInfo(ctx context.Context, socialInfo model.SocialInfo) (*model.SocialInfo, error) {
	return b.db.UpsertSocialInfo(ctx, socialInfo)
}
