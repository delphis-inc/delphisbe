package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nedrocks/delphisbe/graph/model"
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

	// res, err := dbSvc.Query(&dynamodb.QueryInput{
	// 	TableName: aws.String(conf.DBConfig.TablesConfig.Users.TableName),
	// })
	res, err := dbSvc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(conf.DBConfig.TablesConfig.Users.TableName),
	})

	if err != nil {
		logrus.WithError(err).Errorf("Failed to query dynamo table")
		return
	}

	userObjs := []model.User{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &userObjs)

	if err != nil {
		logrus.WithError(err).Errorf("Failed to unmarshal values")
		return
	}

	fmt.Printf("All Users: %+v\n", userObjs)

	res, err = dbSvc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(conf.DBConfig.TablesConfig.UserProfiles.TableName),
	})

	if err != nil {
		logrus.WithError(err).Errorf("Failed to query dynamo table")
		return
	}

	userProfileObjs := []model.UserProfile{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &userProfileObjs)

	if err != nil {
		logrus.WithError(err).Errorf("Failed to unmarshal values")
		return
	}

	fmt.Printf("All profiles: %+v\n", userProfileObjs)
}
