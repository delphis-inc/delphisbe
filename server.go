package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials/endpointcreds"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/vektah/gqlparser/v2/formatter"

	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/nedrocks/delphisbe/internal/secrets"

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
	config.AddConfigDirectory("/var/delphis/config")
	conf, err := config.ReadConfig()
	if err != nil {
		logrus.WithError(err).Errorf("Error loading config file")
		return
	}

	awsConfig := aws.NewConfig().WithRegion(conf.AWS.Region).WithCredentialsChainVerboseErrors(true)
	var awsSession *session.Session
	if conf.AWS.UseCredentials {
		awsConfig = awsConfig.WithCredentials(credentials.NewStaticCredentials(
			conf.AWS.Credentials.ID, conf.AWS.Credentials.Secret, conf.AWS.Credentials.Token))
	} else if conf.AWS.IsFargate {
		if ECSCredentialsURI, exists := os.LookupEnv("AWS_CONTAINER_CREDENTIALS_RELATIVE_URI"); exists {
			endpoint := fmt.Sprintf("http://169.254.170.2%s", ECSCredentialsURI)
			awsSession = session.New(awsConfig)
			providerClient := endpointcreds.NewProviderClient(*awsSession.Config, awsSession.Handlers, endpoint)
			creds := credentials.NewCredentials(providerClient)
			awsConfig = awsConfig.WithCredentials(creds)
		}
	}
	awsSession = session.Must(session.NewSession(awsConfig))

	secretManager := secrets.NewSecretsManager(awsConfig, awsSession)
	secrets, err := secretManager.GetSecrets()

	if err == nil {
		for k, v := range secrets {
			os.Setenv(k, v)
		}
		conf.ReadEnvAndUpdate()
	}

	delphisBackend := backend.NewDelphisBackend(*conf, awsSession)

	generatedSchema := generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolver.Resolver{
			DAOManager: delphisBackend,
		},
	})

	//wrappedSchema := introspection.WrapSchema(generatedSchema.Schema())
	b := bytes.Buffer{}
	f := formatter.NewFormatter(&b)
	f.FormatSchema(generatedSchema.Schema())

	srv := handler.NewDefaultServer(generatedSchema)

	http.Handle("/", allowCors(healthCheck()))
	http.Handle("/graphiql", allowCors(playground.Handler("GraphQL playground", "/query")))
	http.Handle("/query", allowCors(authMiddleware(*conf, delphisBackend, srv)))
	config := &oauth1.Config{
		ConsumerKey:    conf.Twitter.ConsumerKey,
		ConsumerSecret: conf.Twitter.ConsumerSecret,
		CallbackURL:    conf.Twitter.Callback,
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}

	http.Handle("/twitter/login", twitter.LoginHandler(config, nil))
	http.Handle("/twitter/callback", twitter.CallbackHandler(config, successfulLogin(*conf, delphisBackend), nil))
	http.Handle("/health", healthCheck())
	log.Printf("connect on port %s for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// TODO: This is quite hacky but fulfills our purposes for now.
func authMiddleware(conf config.Config, delphisBackend backend.DelphisBackend, next http.Handler) http.Handler {
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
						Domain:   conf.Auth.Domain,
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

func allowCors(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		referrer := req.Header.Get("Referer")
		parsedURL, err := url.Parse(referrer)
		if err == nil {
			parts := strings.Split(parsedURL.Host, ":")
			if len(parts) > 0 && strings.HasSuffix(parts[0], "delphishq.com") {
				w.Header().Add("Access-Control-Allow-Origin", fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host))
				w.Header().Add("Access-Control-Allow-Headers", "Host, Accept-Encoding, Accept, Referer, Sec-Fetch-Dest, User-Agent, Content-Type, Content-Length")
				w.Header().Add("Access-Control-Allow-Credentials", "true")
			}
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func healthCheck() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "success")
	}
	return http.HandlerFunc(fn)
}

func successfulLogin(conf config.Config, delphisBackend backend.DelphisBackend) http.Handler {
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
			Domain:   conf.Auth.Domain,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
			MaxAge:   int(30 * 24 * time.Hour / time.Second),
			HttpOnly: true,
		})

		redirectURL, err := url.Parse(conf.Twitter.Redirect)
		if err != nil {
			logrus.WithError(err).Fatalf("Failed parsing the redirect URI so failing the app")
		}
		// TODO: Only want to do this if authing via the app. Given that will
		// be the majority use case, just doing this for everyone... for now.
		query := redirectURL.Query()
		query.Set("dc", authToken.TokenString)
		redirectURL.RawQuery = query.Encode()
		http.Redirect(w, req, redirectURL.String(), 302)
	}
	return http.HandlerFunc(fn)
}
