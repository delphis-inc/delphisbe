package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/nedrocks/delphisbe/internal/auth"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	gologinOauth1 "github.com/dghubble/gologin/oauth1"
	"github.com/dghubble/gologin/twitter"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
	"github.com/nedrocks/delphisbe/graph/generated"
	"github.com/nedrocks/delphisbe/graph/resolver"
	"github.com/nedrocks/delphisbe/internal/backend"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/sirupsen/logrus"
)

const (
	defaultPort = "8080"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	config.AddConfigDirectory("./config")
	conf, err := config.ReadConfig()
	if err != nil {
		logrus.WithError(err).Errorf("Error loading config file")
		return
	}

	delphisBackend := backend.NewDelphisBackend(*conf)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolver.Resolver{
			DAOManager: delphisBackend,
		}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", authMiddleware(delphisBackend, srv))
	config := &oauth1.Config{
		ConsumerKey:    conf.Twitter.ConsumerKey,
		ConsumerSecret: conf.Twitter.ConsumerSecret,
		CallbackURL:    "http://local.delphishq.com:8080/twitter/callback",
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}
	http.Handle("/twitter/login", twitter.LoginHandler(config, nil))
	http.Handle("/twitter/callback", twitter.CallbackHandler(config, successfulLogin(delphisBackend), nil))
	log.Printf("connect to http://local.delphishq.com:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// TODO: This is quite hacky but fulfills our purposes for now.
func authMiddleware(delphisBackend backend.DelphisBackend, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var accessTokenString string
		req := r
		isCookie := false
		accessTokenStringArr := r.URL.Query()["access_token"]
		if len(accessTokenStringArr) > 0 {
			accessTokenString = accessTokenStringArr[0]
		} else {
			accessTokenCookie, err := r.Cookie("delphis_access_token")
			if accessTokenCookie != nil && err == nil {
				accessTokenString = accessTokenCookie.Value
				isCookie = true
			}
		}
		if accessTokenString != "" {
			authedUser, err := delphisBackend.ValidateAccessToken(r.Context(), accessTokenString)
			if err != nil || authedUser == nil {
				if isCookie {
					http.SetCookie(w, &http.Cookie{
						Name:     "delphis_access_token",
						Value:    "",
						Domain:   "local.delphishq.com",
						Path:     "/",
						MaxAge:   0,
						HttpOnly: true,
						SameSite: http.SameSiteStrictMode,
					})
				}
			} else {
				ctx := auth.WithAuthedUser(r.Context(), authedUser)
				req = r.WithContext(ctx)
			}
		}
		next.ServeHTTP(w, req)
	})
}

func successfulLogin(delphisBackend backend.DelphisBackend) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		twitterUser, err := twitter.UserFromContext(ctx)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to extract twitter user from context")
			return
		}
		accessToken, accessTokenSecret, err := gologinOauth1.AccessTokenFromContext(ctx)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to extract oath tokens from request context")
			return
		}

		userObj, err := delphisBackend.GetOrCreateUser(ctx, backend.LoginWithTwitterInput{
			User:              twitterUser,
			AccessToken:       accessToken,
			AccessTokenSecret: accessTokenSecret,
		})
		if err != nil {
			logrus.WithError(err).Errorf("Got an error creating a user")
			return
		}

		// At this point we hae a successful login
		authToken, err := delphisBackend.NewAccessToken(req.Context(), userObj.ID)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to create access token")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "delphis_access_token",
			Value:    authToken.TokenString,
			Domain:   "local.delphishq.com",
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
			MaxAge:   int(30 * 24 * time.Hour / time.Second),
			HttpOnly: true,
		})
	}
	return http.HandlerFunc(fn)
}