package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/delphis-inc/delphisbe/internal/config"
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

	result, err := dbSvc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		logrus.WithError(err).Errorf("Could not list tables")
		return
	}

	logrus.Println("Tables:")
	for _, table := range result.TableNames {
		logrus.Println(*table)
	}

	res, err := dbSvc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(conf.DBConfig.TablesConfig.Users.TableName),
	})

	logrus.Println("User table:")
	logrus.Println(res.String())
}
