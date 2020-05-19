package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials/endpointcreds"
	"github.com/gorilla/websocket"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/vektah/gqlparser/v2/formatter"

	"github.com/nedrocks/delphisbe/internal/auth"
	"github.com/nedrocks/delphisbe/internal/secrets"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
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
	logrus.Debugf("Starting")

	rand.Seed(time.Now().Unix())

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
	logrus.Debugf("Got config from file")

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
	logrus.Debugf("Got creds from remote")
	awsSession = session.Must(session.NewSession(awsConfig))

	secretManager := secrets.NewSecretsManager(awsConfig, awsSession)
	secrets, err := secretManager.GetSecrets()
	logrus.Debugf("Got secrets")
	if err == nil {
		for k, v := range secrets {
			os.Setenv(k, v)
		}
		conf.ReadEnvAndUpdate()
	}

	logrus.Debugf("about to create backend")
	delphisBackend := backend.NewDelphisBackend(*conf, awsSession)
	logrus.Debugf("Created backend")

	generatedSchema := generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolver.Resolver{
			DAOManager: delphisBackend,
		},
	})

	b := bytes.Buffer{}
	f := formatter.NewFormatter(&b)
	f.FormatSchema(generatedSchema.Schema())

	srv := handler.NewDefaultServer(generatedSchema)

	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

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
	http.Handle("/upload_image", allowCors(uploadImage(delphisBackend)))
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
		// Check Headers
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			accessTokenString = strings.Split(authHeader, " ")[1]
		}
		// Check Query String (overrides header)
		accessTokenStringArr := r.URL.Query()["access_token"]
		if len(accessTokenStringArr) > 0 {
			accessTokenString = accessTokenStringArr[0]
		}
		// Check cookie (overrides header and query string)
		accessTokenCookie, err := r.Cookie("delphis_access_token")
		if accessTokenCookie != nil && err == nil {
			accessTokenString = accessTokenCookie.Value
			isCookie = true
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

func uploadImage(delphisBackend backend.DelphisBackend) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			logrus.WithError(errors.New("non-POST request was sent to uploadImage"))
			w.WriteHeader(http.StatusBadRequest)
		}

		// Limit to 10MB
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			logrus.WithError(err).Error("uploaded image was over 10MB")
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("File was over 10MB")); err != nil {
				return
			}
			return
		}

		// Retrieve image file
		file, header, err := r.FormFile("image")
		if err != nil {
			logrus.WithError(err).Error("failed getting image from form file")
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(fmt.Sprintf("500 - Something bad happened!"))); err != nil {
				return
			}
			return
		}

		// Check for an empty file
		if header.Size == 0 {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("400 - File was empty!")); err != nil {
				return
			}
			return
		}

		// Get file extension
		// The client should always send a properly formed media file
		ext := filepath.Ext(header.Filename)
		if len(ext) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("400 - No extension on media file")); err != nil {
				return
			}
			return
		}

		// Upload image
		mediaID, mimeType, err := delphisBackend.UploadMedia(r.Context(), ext, file)
		if err != nil {
			logrus.WithError(err).Error("failed to upload media")
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(fmt.Sprintf("500 - Something bad happened!"))); err != nil {
				return
			}
			return
		}

		w.Header().Set("Content-type", " application/json")

		// TODO: Create a struct if we decide to return more than the ID
		resp := map[string]string{"media_id": mediaID, "media_type": mimeType}
		json.NewEncoder(w).Encode(resp)
	}

	return http.HandlerFunc(fn)
}
