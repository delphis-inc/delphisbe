package main

import (
	"log"
	"net/http"
	"os"

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

const defaultPort = "8080"

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
	http.Handle("/query", srv)
	config := &oauth1.Config{
		ConsumerKey:    conf.Twitter.ConsumerKey,
		ConsumerSecret: conf.Twitter.ConsumerSecret,
		CallbackURL:    "http://localhost:8080/twitter/callback",
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}
	http.Handle("/twitter/login", twitter.LoginHandler(config, nil))
	http.Handle("/twitter/callback", twitter.CallbackHandler(config, successfulLogin(delphisBackend), nil))
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func successfulLogin(delphisBackend backend.DelphisBackend) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		twitterUser, err := twitter.UserFromContext(ctx)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to extract twitter user from context")
		}
		accessToken, accessTokenSecret, err := gologinOauth1.AccessTokenFromContext(ctx)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to extract oath tokens from request context")
		}

		_, err = delphisBackend.GetOrCreateUser(ctx, backend.LoginWithTwitterInput{
			User:              twitterUser,
			AccessToken:       accessToken,
			AccessTokenSecret: accessTokenSecret,
		})
		if err != nil {
			logrus.WithError(err).Errorf("Got an error creating a user")
		}
	}
	return http.HandlerFunc(fn)
}
