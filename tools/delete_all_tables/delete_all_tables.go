package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nedrocks/delphisbe/internal/config"
	"github.com/sirupsen/logrus"
)

func main() {
	config.AddConfigDirectory("./config")
	conf, err := config.ReadConfig()
	if err != nil {
		logrus.WithError(err).Errorf("Error loading config file")
		return
	}
	creds := credentials.NewStaticCredentials("fakeMyKeyId", "fakeSecretAccessKey", "")
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("us-west-2"),
		Endpoint:    aws.String("http://localhost:7998")})

	if err != nil {
		logrus.Println(err)
		return
	}
	dbSvc := dynamodb.New(sess)

	for _, tableName := range []string{
		conf.DBConfig.TablesConfig.Discussions.TableName,
		conf.DBConfig.TablesConfig.Participants.TableName,
		conf.DBConfig.TablesConfig.PostBookmarks.TableName,
		conf.DBConfig.TablesConfig.Posts.TableName,
		conf.DBConfig.TablesConfig.UserProfiles.TableName,
		conf.DBConfig.TablesConfig.Users.TableName,
		conf.DBConfig.TablesConfig.Viewers.TableName,
	} {
		_, err = dbSvc.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})

		if err != nil {
			logrus.WithError(err).Errorf("Failed deleting table: %s.", tableName)
		}
	}
}
