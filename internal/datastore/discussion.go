package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error) {
	discussions, err := d.GetDiscussionsByIDs(ctx, []string{id})
	if err != nil {
		return nil, err
	}

	return discussions[id], nil
}

func (d *db) GetDiscussionsByIDs(ctx context.Context, ids []string) (map[string]*model.Discussion, error) {
	logrus.Debug("GetDiscussionsByIDs: Dynamo BatchGetItem")
	keys := make([]map[string]*dynamodb.AttributeValue, 0)
	for _, id := range ids {
		keys = append(keys, map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		})
	}
	res, err := d.dynamo.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			d.dbConfig.Discussions.TableName: {
				Keys: keys,
			},
		},
	})

	if err != nil {
		logrus.WithError(err).Errorf("GetDiscussionsByIDs: Failed to retrive discussions with ids: %+v", ids)
		return nil, err
	}

	discussionMap := map[string]*model.Discussion{}
	for _, id := range ids {
		discussionMap[id] = nil
	}
	elems := res.Responses[d.dbConfig.Discussions.TableName]
	for _, elem := range elems {
		discussionObj := model.Discussion{}
		err = dynamodbattribute.UnmarshalMap(elem, &discussionObj)
		if err != nil {
			logrus.WithError(err).Warnf("GetDiscussionsByIDs: Failed to unmarshal discussion object: %+v", elem)
			continue
		}

		discussionMap[discussionObj.ID] = &discussionObj
	}

	return discussionMap, nil
}

func (d *db) ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error) {
	logrus.Debug("ListDiscussions::Dynamo Scan")
	res, err := d.dynamo.Scan(&dynamodb.ScanInput{
		TableName: aws.String(d.dbConfig.Discussions.TableName),
	})

	if err != nil {
		logrus.WithError(err).Errorf("ListDiscussions: Failed listing discussions")
		return nil, err
	}

	if res.Count == nil || res.Items == nil {
		logrus.Errorf("ListDiscussions: Returned item set is nil")
	}

	ids := make([]string, 0)
	edges := make([]*model.DiscussionsEdge, 0)
	for _, elem := range res.Items {
		discussionObj := model.Discussion{}
		err := dynamodbattribute.UnmarshalMap(elem, &discussionObj)
		if err != nil {
			logrus.WithError(err).Warnf("ListDiscussion: Failed unmarshaling discussion: %+v", elem)
			continue
		}
		edges = append(edges, &model.DiscussionsEdge{
			Node: &discussionObj,
		})
		ids = append(ids, discussionObj.ID)
	}

	return &model.DiscussionsConnection{
		IDs:   ids,
		Edges: edges,
	}, nil
}

func (d *db) PutDiscussion(ctx context.Context, discussion model.Discussion) (*model.Discussion, error) {
	logrus.Debug("PutDiscussion::Dynamo PutItem")
	av, err := d.marshalMap(discussion)
	if err != nil {
		logrus.WithError(err).Errorf("PutDiscussion: Failed to marshal discussion object: %+v", discussion)
		return nil, err
	}
	_, err = d.dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.dbConfig.Discussions.TableName),
		Item:      av,
	})

	if err != nil {
		logrus.WithError(err).Errorf("PutDiscussion: Failed to put discussion object: %+v", av)
		return nil, err
	}
	return &discussion, nil
}
