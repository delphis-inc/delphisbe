package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/delphis-inc/delphisbe/graph/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/endpointcreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/delphis-inc/delphisbe/internal/backend"
	"github.com/delphis-inc/delphisbe/internal/config"
	"github.com/delphis-inc/delphisbe/internal/secrets"
	"github.com/sirupsen/logrus"
)

const (
	defaultPort = "8080"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
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

	doWork(ctx, delphisBackend)

}

func doWork(ctx context.Context, delphisBackend backend.DelphisBackend) {
	// Fetch discussions
	connection, err := delphisBackend.ListDiscussions(ctx)
	if err != nil {
		panic(err)
	}
	discussions := make([]*model.Discussion, 0)
	for i, edge := range connection.Edges {
		if edge != nil {
			discussions = append(discussions, connection.Edges[i].Node)
		}
	}

	// Iterate over discussions. Check if the concierge user has a participant, if not add one
	for _, disc := range discussions {
		if _, err := delphisBackend.UpsertInviteLinksByDiscussionID(ctx, disc.ID); err != nil {
			logrus.WithError(err).Errorf("failed to upsert invite links for discussion: %v\n", disc.ID)
			panic(err)
		}
	}
}
