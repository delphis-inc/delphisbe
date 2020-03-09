package datastore

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/config"
)

func Test_MarshalPost(t *testing.T) {
	type args struct {
		post model.Post
	}

	havePostObj := model.Post{
		ID:            "11111",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DiscussionID:  "22222",
		ParticipantID: 99999,
		PostContentID: "33333",
		Participant:   &model.Participant{},
		PostContent: model.PostContent{
			ID:      "33333",
			Content: "Lorem ipsum dolar amet",
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
				post: havePostObj,
			},
			want: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(havePostObj.ID),
				},
				"CreatedAt": {
					S: aws.String(havePostObj.CreatedAt.Format(time.RFC3339Nano)),
				},
				"UpdatedAt": {
					S: aws.String(havePostObj.UpdatedAt.Format(time.RFC3339Nano)),
				},
				"DeletedAt": {
					NULL: aws.Bool(true),
				},
				"DeletedReasonCode": {
					NULL: aws.Bool(true),
				},
				"DiscussionID": {
					S: aws.String(havePostObj.DiscussionID),
				},
				"ParticipantID": {
					N: aws.String(strconv.Itoa(havePostObj.ParticipantID)),
				},
				"PostContentID": {
					S: aws.String(havePostObj.PostContentID),
				},
				"PostContent": {
					M: map[string]*dynamodb.AttributeValue{
						"ID": {
							S: aws.String(havePostObj.PostContent.ID),
						},
						"Content": {
							S: aws.String(havePostObj.PostContent.Content),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshaled, err := datastoreObj.marshalMap(tt.args.post)
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
