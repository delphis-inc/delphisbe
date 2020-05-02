package main

import (
    "bufio"
    "os"
    "fmt"
    "encoding/csv"

    "github.com/aws/aws-sdk-go/aws/credentials/endpointcreds"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"

    "github.com/nedrocks/delphisbe/internal/secrets"

    "github.com/nedrocks/delphisbe/internal/backend"
    "github.com/nedrocks/delphisbe/internal/config"
    "github.com/sirupsen/logrus"
    flair_templates "./csv_loaders"
)

func main() {
    logrus.SetLevel(logrus.DebugLevel)
    logrus.Debugf("Starting")

	// Read in args
    if len(os.Args) != 2 {
        logrus.Fatal("Usage: load_from_csv.go object csv_filename")
    }
    object := os.Args[1]
    filename := os.Args[2]

    // Load delphisBackend
    config.AddConfigDirectory("./config")
    config.AddConfigDirectory("/var/delphis/config")
    conf, err := config.ReadConfig()
    if err != nil {
        logrus.WithError(err).Fatal("Error loading config file")
    }
    logrus.Debug("Got config from file")

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
    logrus.Debug("Got secrets")
    if err == nil {
        for k, v := range secrets {
            os.Setenv(k, v)
        }
        conf.ReadEnvAndUpdate()
    }

    logrus.Debug("about to create backend")
    delphisBackend := backend.NewDelphisBackend(*conf, awsSession)
    logrus.Debug("Created backend")

    // Read csv
    csvFile, _ := os.Open(filename)
    reader := csv.NewReader(bufio.NewReader(csvFile))

    logrus.Debug(reader)
    logrus.Debug(delphisBackend)
    create_from_csv(delphisBackend, reader)
}
