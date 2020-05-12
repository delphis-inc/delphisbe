package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *db) GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error) {
	logrus.Debug("GetDiscussionByID::SQL Query")
	discussions, err := d.GetDiscussionsByIDs(ctx, []string{id})
	if err != nil {
		return nil, err
	}

	return discussions[id], nil
}

func (d *db) GetDiscussionsByIDs(ctx context.Context, ids []string) (map[string]*model.Discussion, error) {
	logrus.Debug("GetDiscussionsByIDs::SQL Query")
	discussions := []model.Discussion{}
	if err := d.sql.Where(ids).Preload("Moderator").Find(&discussions).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// This is a not found situation with multiple ids. Not sure what to do here..?
		} else {
			logrus.WithError(err).Errorf("GetDiscussionsByIDs::Failed to get discussions by IDs")
			return nil, err
		}
	}
	retVal := map[string]*model.Discussion{}
	for _, id := range ids {
		retVal[id] = nil
	}
	for _, disc := range discussions {
		retVal[disc.ID] = &disc
	}
	return retVal, nil
}

func (d *db) GetDiscussionByModeratorID(ctx context.Context, moderatorID string) (*model.Discussion, error) {
	logrus.Debugf("GetDiscussionByModeratorID::SQL Query")
	discussion := model.Discussion{}
	moderator := model.Moderator{}
	if err := d.sql.Preload("Moderator").First(&moderator, model.Moderator{ID: moderatorID}).Related(&discussion).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetDiscussionByModeratorID::Failed getting discussion by moderator ID")
		return nil, err
	}

	return &discussion, nil
}

func (d *db) ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error) {
	//TODO: this should take in paging params and return based on those.
	logrus.Debugf("ListDiscussions::SQL Query")

	discussions := []model.Discussion{}
	if err := d.sql.Preload("Moderator").Find(&discussions).Error; err != nil {
		logrus.WithError(err).Errorf("ListDiscussions::Failed to list discussions")
		return nil, err
	}

	ids := make([]string, 0)
	edges := make([]*model.DiscussionsEdge, 0)
	for i := range discussions {
		discussionObj := &discussions[i]
		edges = append(edges, &model.DiscussionsEdge{
			Node: discussionObj,
		})
		ids = append(ids, discussionObj.ID)
	}

	return &model.DiscussionsConnection{
		Edges: edges,
		IDs:   ids,
	}, nil
}

func (d *db) UpsertDiscussion(ctx context.Context, discussion model.Discussion) (*model.Discussion, error) {
	logrus.Debug("UpsertDiscussion::SQL Create")
	found := model.Discussion{}
	if err := d.sql.First(&found, model.Discussion{ID: discussion.ID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := d.sql.Create(&discussion).First(&found, model.Discussion{ID: discussion.ID}).Error; err != nil {
				logrus.WithError(err).Errorf("UpsertDiscussion::Failed to create new object")
				return nil, err
			}
		} else {
			logrus.WithError(err).Errorf("UpsertDiscussion::Failed checking for Discussion object")
			return nil, err
		}
	} else {
		if err := d.sql.Preload("Moderator").Model(&discussion).Updates(model.Discussion{
			Title:         discussion.Title,
			AnonymityType: discussion.AnonymityType,
		}).First(&found).Error; err != nil {
			logrus.WithError(err).Errorf("UpsertDiscussion::Failed updating disucssion object")
			return nil, err
		}
	}
	return &found, nil
}

//////////
//Dynamo functions
//////////
func (d *db) PutDiscussionDynamo(ctx context.Context, discussion model.Discussion) (*model.Discussion, error) {
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

func (d *db) ListDiscussionsDynamo(ctx context.Context) (*model.DiscussionsConnection, error) {
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

func (d *db) GetDiscussionsByIDsDynamo(ctx context.Context, ids []string) (map[string]*model.Discussion, error) {
	logrus.Debug("GetDiscussionsByIDs: Dynamo BatchGetItem")
	keys := make([]map[string]*dynamodb.AttributeValue, 0)
	for _, id := range ids {
		keys = append(keys, map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		})
	}
	// NOTE: This has to be BatchGet because query will not work unless we specify partition key.
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
