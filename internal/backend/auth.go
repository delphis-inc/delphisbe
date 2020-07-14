package backend

import (
	"context"

	"github.com/delphis-inc/delphisbe/internal/auth"
)

func (b *delphisBackend) NewAccessToken(ctx context.Context, userID string) (*auth.DelphisAccessToken, error) {
	return b.auth.NewAccessToken(userID)
}

func (b *delphisBackend) ValidateAccessToken(ctx context.Context, token string) (*auth.DelphisAuthedUser, error) {
	return b.auth.ValidateAccessToken(ctx, token)
}

func (b *delphisBackend) ValidateRefreshToken(ctx context.Context, token string) (*auth.DelphisRefreshTokenUser, error) {
	return b.auth.ValidateRefreshToken(ctx, token)
}
