package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials/endpointcreds"
	"github.com/delphis-inc/delphisbe/internal/worker"
	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/vektah/gqlparser/v2/formatter"

	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/secrets"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/delphis-inc/delphisbe/graph/generated"
	"github.com/delphis-inc/delphisbe/graph/resolver"
	"github.com/delphis-inc/delphisbe/internal/backend"
	"github.com/delphis-inc/delphisbe/internal/config"
	gologinOauth1 "github.com/dghubble/gologin/oauth1"
	"github.com/dghubble/gologin/twitter"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
	"github.com/sirupsen/logrus"
)

const (
	defaultPort = "8080"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	logrus.Debugf("Starting")

	ctx := context.Background()
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

	// Kickoff sqs workers
	dripWorker := worker.NewDripWorker(*conf, delphisBackend, awsSession)
	go func() {
		logrus.Debugf("Kicking off drip sqs")
		for {
			if err := dripWorker.Start(ctx); err != nil {
				logrus.WithError(err).Error("failed to start drip worker")
				time.Sleep(1 * time.Second)
			}
		}
	}()

	// Kickoff cron job
	c := cron.New()
	c.AddFunc("@every 5m", delphisBackend.AutoPostContent)
	c.Start()

	http.Handle("/", allowCors(healthCheck()))
	http.Handle("/graphiql", allowCors(playground.Handler("GraphQL playground", "/query")))
	http.Handle("/query", allowCors(authMiddleware(*conf, delphisBackend, srv)))
	config := &oauth1.Config{
		ConsumerKey:    conf.Twitter.ConsumerKey,
		ConsumerSecret: conf.Twitter.ConsumerSecret,
		CallbackURL:    conf.Twitter.Callback,
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}

	http.Handle("/apple/authLogin", appleAuthLogin(conf, delphisBackend))
	http.Handle("/twitter/login", twitter.LoginHandler(config, nil))
	http.Handle("/twitter/callback", twitter.CallbackHandler(config, successfulLogin(*conf, delphisBackend), nil))
	http.Handle("/upload_image", allowCors(uploadImage(delphisBackend)))
	http.Handle("/health", healthCheck())
	log.Printf("connect on port %s for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func appleAuthLogin(conf *config.Config, delphisBackend backend.DelphisBackend) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.Method != "POST" {
			logrus.WithError(errors.New("non-POST request was sent to appleAuthLogin"))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := r.ParseForm()
		if err != nil {
			logrus.WithError(err).Errorf("Failed to parse request form")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		firstName := r.FormValue("fn")
		lastName := r.FormValue("ln")
		email := r.FormValue("e")
		code := r.FormValue("c")

		if email == "" || code == "" {
			logrus.Errorf("Failed to retrieve email or code while authing")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Now we need to validate the code.
		clientSecretStr, err := auth.GenerateAppleClientSecret(ctx, conf)
		if err != nil || clientSecretStr == nil {
			logrus.WithError(err).Errorf("Did not generate client secret properly.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hc := http.Client{}
		form := url.Values{}
		form.Add("client_id", conf.AppleAuthConfig.ClientID)
		form.Add("client_secret", *clientSecretStr)
		form.Add("code", code)
		form.Add("grant_type", "authorization_code")
		postReq, _ := http.NewRequest("POST", "https://appleid.apple.com/auth/token", strings.NewReader(form.Encode()))

		postReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := hc.Do(postReq)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to make auth/token request to apple")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if resp.StatusCode != 200 {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			logrus.Infof("Failed to exchange code with apple. Received status code: %d. Response was: %s", resp.StatusCode, string(bodyBytes))
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Otherwise this has parsed and we can create the user!
		user, err := delphisBackend.GetOrCreateAppleUser(ctx, backend.LoginWithAppleInput{
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
		})
		if err != nil {
			logrus.WithError(err).Errorf("Failed to create user for apple login")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		authToken, err := delphisBackend.NewAccessToken(ctx, user.ID)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to create access token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		response := map[string]string{"delphis_access_token": authToken.TokenString}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			logrus.WithError(err).Errorf("Failed to encode access token response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(fn)
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
			return
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

		// Upload image
		mediaID, mimeType, err := delphisBackend.UploadMedia(r.Context(), file)
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
