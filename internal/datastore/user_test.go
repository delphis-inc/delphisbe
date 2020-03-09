package datastore

import (
	"reflect"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
)

func Test_MarshalUser(t *testing.T) {
	type args struct {
		user model.User
	}

	haveUserObj := model.User{
		ID:            "11111",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		UserProfileID: "22222",
		UserProfile:   &model.UserProfile{},
		DiscussionParticipants: &model.DiscussionParticipantKeys{
			Keys: []model.DiscussionParticipantKey{
				{
					DiscussionID:  "33333",
					ParticipantID: 0,
				},
			},
		},
		DiscussionViewers: &model.DiscussionViewerKeys{
			Keys: []model.DiscussionViewerKey{
				{
					DiscussionID: "33333",
					ViewerID:     "44444",
				},
			},
		},
		Participants: []*model.Participant{},
		Viewers:      []*model.Viewer{},
	}

	datastoreObj := NewDatastore(config.DBConfig{})

	tests := []struct {
		name string
		args args
		want map[string]*dynamodb.AttributeValue
	}{
		{
			name: "fully filled object",
			args: args{
				user: haveUserObj,
			},
			want: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(haveUserObj.ID),
				},
				"CreatedAt": {
					S: aws.String(haveUserObj.CreatedAt.Format(time.RFC3339Nano)),
				},
				"UpdatedAt": {
					S: aws.String(haveUserObj.UpdatedAt.Format(time.RFC3339Nano)),
				},
				"DeletedAt": {
					NULL: aws.Bool(true),
				},
				"UserProfileID": {
					S: aws.String(haveUserObj.UserProfileID),
				},
				"DiscussionParticipants": {
					SS: []*string{aws.String(haveUserObj.DiscussionParticipants.Keys[0].String())},
				},
				"DiscussionViewers": {
					SS: []*string{aws.String(haveUserObj.DiscussionViewers.Keys[0].String())},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshaled, err := datastoreObj.marshalMap(tt.args.user)
			if err != nil {
				t.Errorf("Caught an error marshaling: %+v", err)
				return
			}
			if !reflect.DeepEqual(marshaled, tt.want) {
				t.Errorf("These objects did not match. Got: %+v\n\n Want: %+v", marshaled, tt.want)
			}
		})
	}
}
