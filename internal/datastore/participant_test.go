package datastore

import (
	"testing"
)

func Test_MarshalParticipant(t *testing.T) {
	// type args struct {
	// 	participant model.Participant
	// }

	// haveParticipantObj := model.Participant{
	// 	ParticipantID:                     11111,
	// 	CreatedAt:                         time.Now(),
	// 	UpdatedAt:                         time.Now(),
	// 	DiscussionID:                      "12345",
	// 	ViewerID:                          "54321",
	// 	DiscussionNotificationPreferences: model.ParticipantNotificationPreferences{},
	// 	Viewer:                            &model.Viewer{},
	// 	Discussion:                        &model.Discussion{},
	// 	Posts:                             &model.PostsConnection{},
	// 	UserID:                            "22222",
	// 	User:                              &model.User{},
	// }

	// datastoreObj := NewDatastore(config.DBConfig{})

	// tests := []struct {
	// 	name string
	// 	args args
	// 	want map[string]*dynamodb.AttributeValue
	// }{
	// 	{
	// 		name: "fully filled object",
	// 		args: args{
	// 			participant: haveParticipantObj,
	// 		},
	// 		want: map[string]*dynamodb.AttributeValue{
	// 			"ParticipantID": {
	// 				N: aws.String(strconv.Itoa(haveParticipantObj.ParticipantID)),
	// 			},
	// 			"CreatedAt": {
	// 				S: aws.String(haveParticipantObj.CreatedAt.Format(time.RFC3339Nano)),
	// 			},
	// 			"UpdatedAt": {
	// 				S: aws.String(haveParticipantObj.UpdatedAt.Format(time.RFC3339Nano)),
	// 			},
	// 			"DeletedAt": {
	// 				NULL: aws.Bool(true),
	// 			},
	// 			"DiscussionID": {
	// 				S: aws.String(haveParticipantObj.DiscussionID),
	// 			},
	// 			"ViewerID": {
	// 				S: aws.String(haveParticipantObj.ViewerID),
	// 			},
	// 			"DiscussionNotificationPreferences": {
	// 				M: map[string]*dynamodb.AttributeValue{
	// 					"ID": {
	// 						NULL: aws.Bool(true),
	// 					},
	// 				},
	// 			},
	// 			"UserID": {
	// 				S: aws.String(haveParticipantObj.UserID),
	// 			},
	// 		},
	// 	},
	// }

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		marshaled, err := datastoreObj.marshalMap(tt.args.participant)
	// 		if err != nil {
	// 			t.Errorf("Caught an error marshaling: %+v", err)
	// 			return
	// 		}
	// 		if !reflect.DeepEqual(marshaled, tt.want) {
	// 			t.Errorf("These objects did not match. Got: %+v\n\n Want: %+v", marshaled, tt.want)
	// 		}
	// 	})
	// }
}
