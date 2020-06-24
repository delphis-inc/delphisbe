package datastore

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) GetParticipantByID(ctx context.Context, id string) (*model.Participant, error) {
	logrus.Debugf("GetParticipantByID::SQL Query")
	participants, err := d.GetParticipantsByIDs(ctx, []string{id})
	if err != nil {
		return nil, err
	}

	return participants[id], nil
}

func (d *delphisDB) GetParticipantsByIDs(ctx context.Context, ids []string) (map[string]*model.Participant, error) {
	logrus.Debugf("GetParticipantsByIDs::SQL Query")
	participants := []model.Participant{}
	if err := d.sql.Where(ids).Find(&participants).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetParticipantByID::Failed to query participant by ID")
		return nil, err
	}

	retVal := map[string]*model.Participant{}
	for _, p := range participants {
		tempVal := p
		retVal[p.ID] = &tempVal
	}

	return retVal, nil
}

func (d *delphisDB) GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error) {
	logrus.Debugf("GetParticipantsByDiscussionID::SQL Query")
	participants := []model.Participant{}
	if err := d.sql.Where(&model.Participant{DiscussionID: &id}).Order("participant_id desc").Find(&participants).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetParticipantsByDiscussionID::Failed to get participants by discussion ID")
		return nil, err
	}
	return participants, nil
}

func (d *delphisDB) GetParticipantsByDiscussionIDUserID(ctx context.Context, discussionID string, userID string) ([]model.Participant, error) {
	logrus.Debugf("GetParticipantByDiscussionIDUserID::SQL Query")
	participants := []model.Participant{}
	if err := d.sql.Where(&model.Participant{DiscussionID: &discussionID, UserID: &userID}).Order("participant_id desc").Limit(2).Find(&participants).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetParticipantByDiscussionIDUserID::Failed to get participant by discussion ID and user ID")
		return nil, err
	}
	return participants, nil
}

func (d *delphisDB) UpsertParticipant(ctx context.Context, participant model.Participant) (*model.Participant, error) {
	logrus.Debug("UpsertParticipant::SQL Create")
	found := model.Participant{}
	if err := d.sql.First(&found, model.Participant{ID: participant.ID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := d.sql.Create(&participant).First(&found, model.Participant{ID: participant.ID}).Error; err != nil {
				logrus.WithError(err).Errorf("UpsertParticipant::Failed to put Participant")
				return nil, err
			}
		} else {
			logrus.WithError(err).Errorf("UpsertParticipant::Failed checking for Participant object")
			return nil, err
		}
	} else {
		if err := d.sql.Model(&participant).Updates(model.Participant{
			FlairID:       participant.FlairID,
			IsAnonymous:   participant.IsAnonymous,
			UpdatedAt:     time.Now(),
			GradientColor: participant.GradientColor,
			HasJoined:     participant.HasJoined,
		}).First(&found).Error; err != nil {
			logrus.WithError(err).Errorf("UpsertParticipant::Failed updating Participant object")
			return nil, err
		}
	}
	return &found, nil
}

func (d *delphisDB) AssignFlair(ctx context.Context, participant model.Participant, flairID *string) (*model.Participant, error) {
	logrus.Debug("AssignFlair::SQL Update")
	if err := d.sql.Model(&participant).Update("FlairID", flairID).Error; err != nil {
		logrus.WithError(err).Errorf("AssignFlair::Failed to update")
		return &participant, err
	}
	return &participant, nil
}

func (d *delphisDB) GetTotalParticipantCountByDiscussionID(ctx context.Context, discussionID string) int {
	count := 0
	d.sql.Model(&model.Participant{}).Where(&model.Participant{DiscussionID: &discussionID}).Count(&count)
	return count
}

////////////
//Dynamo functions
////////////

// func (d *db) GetParticipantsByDiscussionIDDynamo(ctx context.Context, id string) ([]model.Participant, error) {
// 	logrus.Debug("GetParticipantsByDiscussionID::Dynamo Query")
// 	res, err := d.dynamo.Query(&dynamodb.QueryInput{
// 		TableName: aws.String(d.dbConfig.Participants.TableName),
// 		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
// 			":did": {
// 				S: aws.String(id),
// 			},
// 		},
// 		KeyConditionExpression: aws.String("DiscussionID = :did"),
// 	})

// 	if err != nil {
// 		logrus.WithError(err).Errorf("GetParticipantsByDiscussionID: Failed to query participants for discussionID: %s", id)
// 		return nil, err
// 	}

// 	participants := make([]model.Participant, 0)
// 	if res != nil {
// 		for _, elem := range res.Items {
// 			participantObj := model.Participant{}
// 			err := dynamodbattribute.UnmarshalMap(elem, &participantObj)
// 			if err != nil {
// 				logrus.WithError(err).Warnf("GetParticipantsByDiscussionID: Failed unmarshaling participant object: %+v", elem)
// 				continue
// 			}
// 			participants = append(participants, participantObj)
// 		}
// 	}
// 	return participants, nil
// }

// func (d *db) GetParticipantsByIDsDynamo(ctx context.Context, discussionParticipantKeys []model.DiscussionParticipantKey) (map[model.DiscussionParticipantKey]*model.Participant, error) {
// 	if len(discussionParticipantKeys) == 0 {
// 		return map[model.DiscussionParticipantKey]*model.Participant{}, nil
// 	}
// 	logrus.Debug("GetParticipantsByIDs::Dynamo BatchGetItem")
// 	// NOTE: Unless we are fetching from the same discussion we need to use BatchGetItem instead
// 	// of Query here.
// 	keys := make([]map[string]*dynamodb.AttributeValue, 0)
// 	for _, dp := range discussionParticipantKeys {
// 		keys = append(keys, map[string]*dynamodb.AttributeValue{
// 			"DiscussionID": {
// 				S: aws.String(dp.DiscussionID),
// 			},
// 			"ParticipantID": {
// 				N: aws.String(strconv.Itoa(dp.ParticipantID)),
// 			},
// 		})
// 	}
// 	res, err := d.dynamo.BatchGetItem(&dynamodb.BatchGetItemInput{
// 		RequestItems: map[string]*dynamodb.KeysAndAttributes{
// 			d.dbConfig.Participants.TableName: {
// 				Keys: keys,
// 			},
// 		},
// 	})

// 	if err != nil {
// 		logrus.WithError(err).Errorf("GetParticipantsByIDs: Failed to query participants for keys: %+v", keys)
// 		return nil, err
// 	}

// 	participantMap := map[model.DiscussionParticipantKey]*model.Participant{}
// 	for _, dp := range discussionParticipantKeys {
// 		participantMap[dp] = nil
// 	}
// 	elems := res.Responses[d.dbConfig.Participants.TableName]
// 	for _, elem := range elems {
// 		participantObj := model.Participant{}
// 		err := dynamodbattribute.UnmarshalMap(elem, &participantObj)
// 		if err != nil {
// 			logrus.WithError(err).Warnf("Failed to unmarshal participant object: %+v", elem)
// 			continue
// 		}

// 		participantMap[participantObj.DiscussionParticipantKey()] = &participantObj
// 	}

// 	return participantMap, nil
// }

// func (d *db) PutParticipantDynamo(ctx context.Context, participant model.Participant) (*model.Participant, error) {
// 	logrus.Debug("PutParticipant::Dynamo PutItem")
// 	av, err := d.marshalMap(participant)
// 	if err != nil {
// 		logrus.WithError(err).Errorf("PutParticipant: Failed to marshal participant object: %+v", participant)
// 		return nil, err
// 	}

// 	_, err = d.dynamo.PutItem(&dynamodb.PutItemInput{
// 		TableName: aws.String(d.dbConfig.Participants.TableName),
// 		Item:      av,
// 	})

// 	if err != nil {
// 		logrus.WithError(err).Errorf("PutParticipant: Failed to put participant object: %+v", av)
// 		return nil, err
// 	}

// 	return &participant, nil
// }
