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

func Test_MarshalDiscussion(t *testing.T) {
	type args struct {
		discussion model.Discussion
	}
	haveDiscussionObj := model.Discussion{
		ID:            "12345",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		AnonymityType: model.AnonymityTypeWeak,
		Posts:         &model.PostsConnection{},
		Participants:  []*model.Participant{},
		Moderator: model.Moderator{
			ID:            "54321",
			DiscussionID:  "12345",
			UserProfileID: "99999",
			UserProfile:   &model.UserProfile{},
			Discussion:    &model.Discussion{},
		},
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
				discussion: haveDiscussionObj,
			},
			want: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(haveDiscussionObj.ID),
				},
				"CreatedAt": {
					S: aws.String(haveDiscussionObj.CreatedAt.Format(time.RFC3339Nano)),
				},
				"UpdatedAt": {
					S: aws.String(haveDiscussionObj.UpdatedAt.Format(time.RFC3339Nano)),
				},
				"DeletedAt": {
					NULL: aws.Bool(true),
				},
				"AnonymityType": {
					S: aws.String(haveDiscussionObj.AnonymityType.String()),
				},
				"Moderator": {
					M: map[string]*dynamodb.AttributeValue{
						"UserProfileID": {
							S: aws.String(haveDiscussionObj.Moderator.UserProfileID),
						},
						"ID": {
							S: aws.String(haveDiscussionObj.Moderator.ID),
						},
						"DiscussionID": {
							S: aws.String(haveDiscussionObj.Moderator.DiscussionID),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshaled, err := datastoreObj.marshalMap(tt.args.discussion)
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
