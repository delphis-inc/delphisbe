package main

import (
	"errors"

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

	err = createDiscussions(dbSvc, conf)
	var perr *dynamodb.ResourceInUseException
	if err != nil {
		if errors.As(err, &perr) {
			// Silently pass on this
			logrus.Debugf("Table Discussions already exists")
		} else {
			logrus.WithError(err).Error("Error when creating discussion table")
		}
	}

	err = createParticipants(dbSvc, conf)
	if err != nil {
		if errors.As(err, &perr) {
			// Silently pass on this
			logrus.Debugf("Table Participants already exists")
		} else {
			logrus.WithError(err).Error("Error when creating participants table")
		}
	}

	err = createPostBookmarks(dbSvc, conf)
	if err != nil {
		if errors.As(err, &perr) {
			// Silently pass on this
			logrus.Debugf("Table Post-Bookmarks already exists")
		} else {
			logrus.WithError(err).Error("Error when creating post-bookmarks table")
		}
	}

	err = createPosts(dbSvc, conf)
	if err != nil {
		if errors.As(err, &perr) {
			// Silently pass on this
			logrus.Debugf("Table Posts already exists")
		} else {
			logrus.WithError(err).Error("Error when creating posts table")
		}
	}

	err = createUsers(dbSvc, conf)
	if err != nil {
		if errors.As(err, &perr) {
			// Silently pass on this
			logrus.Debugf("Table Users already exists")
		} else {
			logrus.WithError(err).Error("Error when creating users table")
		}
	}

	err = createViewers(dbSvc, conf)
	if err != nil {
		if errors.As(err, &perr) {
			// Silently pass on this
			logrus.Debugf("Table Viewers already exists")
		} else {
			logrus.WithError(err).Error("Error when creating viewers table")
		}
	}

	// logrus.Println("Tables:")
	// for _, table := range result.TableNames {
	// 	logrus.Println(*table)
	// }

}

func createDiscussions(dbSvc *dynamodb.DynamoDB, conf *config.Config) error {
	tableName := conf.DBConfig.TablesConfig.Discussions.TableName
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := dbSvc.CreateTable(input)
	return err
}

func createParticipants(dbSvc *dynamodb.DynamoDB, conf *config.Config) error {
	tableName := conf.DBConfig.TablesConfig.Participants.TableName
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("DiscussionID"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("ParticipantID"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("DiscussionID"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("ParticipantID"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := dbSvc.CreateTable(input)
	return err
}

func createPostBookmarks(dbSvc *dynamodb.DynamoDB, conf *config.Config) error {
	tableName := conf.DBConfig.TablesConfig.PostBookmarks.TableName
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := dbSvc.CreateTable(input)
	return err
}

func createPosts(dbSvc *dynamodb.DynamoDB, conf *config.Config) error {
	tableName := conf.DBConfig.TablesConfig.Posts.TableName
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := dbSvc.CreateTable(input)
	return err
}

func createUsers(dbSvc *dynamodb.DynamoDB, conf *config.Config) error {
	tableName := conf.DBConfig.TablesConfig.Users.TableName
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := dbSvc.CreateTable(input)
	return err
}

func createViewers(dbSvc *dynamodb.DynamoDB, conf *config.Config) error {
	tableName := conf.DBConfig.TablesConfig.Viewers.TableName
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String(tableName),
	}

	_, err := dbSvc.CreateTable(input)
	return err
}
