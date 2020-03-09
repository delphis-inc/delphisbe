package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/sirupsen/logrus"
)

const (
	// Will make this shorter in the future but building simple auth token expiry for now.
	authTokenExpiry = 15 * time.Minute
	// 30 days (unused for now)
	refreshTokenExpiry      = 30 * 24 * time.Hour
	issuer                  = "delphishq.com"
	authedUserContextString = "delphisAuthToken"
)

func WithAuthedUser(ctx context.Context, user *DelphisAuthedUser) context.Context {
	return context.WithValue(ctx, authedUserContextString, user)
}

func GetAuthedUser(ctx context.Context) *DelphisAuthedUser {
	val := ctx.Value(authedUserContextString)
	if valAsAuthedUser, ok := val.(*DelphisAuthedUser); ok {
		return valAsAuthedUser
	}
	return nil
}

type DelphisAuth interface {
	NewAccessToken(userID string) (*DelphisAccessToken, error)
	NewRefreshToken(userID string) (*DelphisRefreshToken, error)
	ValidateAccessToken(ctx context.Context, token string) (*DelphisAuthedUser, error)
	ValidateRefreshToken(ctx context.Context, token string) (*DelphisRefreshTokenUser, error)
}

func NewDelphisAuth(config *config.AuthConfig) DelphisAuth {
	return &delphisAuth{
		Config: config,
	}
}

type delphisAuth struct {
	Config *config.AuthConfig
}

func (d *delphisAuth) NewAccessToken(userID string) (*DelphisAccessToken, error) {
	now := time.Now()
	claims := &JWTClaims{
		userID,
		now.Add(authTokenExpiry).Unix(),
		jwt.StandardClaims{
			//ExpiresAt: now.Add(authTokenExpiry).Unix(),
			IssuedAt: now.Unix(),
			Issuer:   issuer,
			Subject:  "at",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(d.Config.HMACSecret))
	if err != nil {
		logrus.WithError(err).Errorf("Failed to sign JWT token")
		return nil, err
	}
	return &DelphisAccessToken{
		Claims:      claims,
		TokenString: signedToken,
	}, nil
}

func (d *delphisAuth) NewRefreshToken(userID string) (*DelphisRefreshToken, error) {
	now := time.Now()
	claims := &JWTClaims{
		userID,
		now.Add(refreshTokenExpiry).Unix(),
		jwt.StandardClaims{
			// Not adding expires at because the library will fail in parsing the token.
			//ExpiresAt: now.Add(refreshTokenExpiry).Unix(),
			IssuedAt: now.Unix(),
			Issuer:   issuer,
			Subject:  "rt",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(d.Config.HMACSecret))
	if err != nil {
		logrus.WithError(err).Errorf("Failed to sign JWT token")
		return nil, err
	}
	return &DelphisRefreshToken{
		Claims:      claims,
		TokenString: signedToken,
	}, nil
}

func (d *delphisAuth) ValidateAccessToken(ctx context.Context, token string) (*DelphisAuthedUser, error) {
	tk, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(d.Config.HMACSecret), nil
	})
	if err != nil {
		logrus.WithError(err).Errorf("Error parsing token")
		return nil, err
	}
	// TODO: Also validate that it doesn't expire longer than expiration expiry in the future.
	if claims, ok := tk.Claims.(*JWTClaims); ok && tk.Valid {
		return &DelphisAuthedUser{
			UserID: claims.UserID,
		}, nil
	}
	err = fmt.Errorf("Failed to validate access token: %s", token)
	logrus.WithError(err).Errorf("Error validating access token")
	return nil, err
}

func (d *delphisAuth) ValidateRefreshToken(ctx context.Context, token string) (*DelphisRefreshTokenUser, error) {
	tk, err := jwt.ParseWithClaims(token, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(d.Config.HMACSecret), nil
	})
	if err != nil {
		logrus.WithError(err).Errorf("Error parsing token")
		return nil, err
	}
	if claims, ok := tk.Claims.(*JWTClaims); ok && tk.Valid {
		return &DelphisRefreshTokenUser{
			UserID: claims.UserID,
		}, nil
	}
	err = fmt.Errorf("Failed to validate refresh token: %s", token)
	logrus.WithError(err).Errorf("Error validating refresh token")
	return nil, err
}

type DelphisAuthedUser struct {
	UserID string
	User   *model.User
}

type DelphisRefreshTokenUser struct {
	UserID string
	User   *model.User
}

type DelphisAccessToken struct {
	Claims      *JWTClaims
	TokenString string
}

type DelphisRefreshToken struct {
	Claims      *JWTClaims
	TokenString string
}

type JWTClaims struct {
	UserID      string
	ReAuthAfter int64
	jwt.StandardClaims
}
