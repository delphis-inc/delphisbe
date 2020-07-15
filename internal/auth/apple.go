package auth

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/url"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/sirupsen/logrus"
)

func AppleAuthorizeURL(redirectURI string, clientID string) string {
	values := url.Values{}
	values.Add("response_type", "code")
	values.Add("redirect_uri", redirectURI)
	values.Add("client_id", clientID)
	values.Add("scope", "name email")
	values.Add("response_mode", "form_post")
	return fmt.Sprintf("https://appleid.apple.com/auth/authorize?%s", values.Encode())
}

func GenerateAppleClientSecret(ctx context.Context, conf *config.Config) (*string, error) {
	privateKeyMayBeBase64 := conf.AppleAuthConfig.PrivateKey
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyMayBeBase64)
	if err != nil {
		privateKeyBytes = []byte(privateKeyMayBeBase64)
	}
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	key, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		logrus.WithError(err).Errorf("Failed decoding apple PKCS8 private key")
	}

	claims := jwt.StandardClaims{
		Issuer:    conf.AppleAuthConfig.TeamID,
		IssuedAt:  time.Now().UTC().Unix(),
		ExpiresAt: time.Now().UTC().Unix() + 24*30*3600,
		Audience:  "https://appleid.apple.com",
		Subject:   conf.AppleAuthConfig.ClientID,
	}
	token := jwt.Token{
		Claims: claims,
		Header: map[string]interface{}{
			"alg": jwt.SigningMethodES256.Alg(),
			"kid": conf.AppleAuthConfig.KeyID,
		},
		Method: jwt.SigningMethodES256,
	}

	signedToken, err := token.SignedString(key)
	if err != nil {
		logrus.WithError(err).Errorf("failed to sign apple JWT token")
		return nil, err
	}

	return &signedToken, nil
}
