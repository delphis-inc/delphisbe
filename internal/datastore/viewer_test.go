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

func Test_MarshalViewer(t *testing.T) {
	type args struct {
		viewer model.Viewer
	}

	lastViewed := time.Now()

	haveViewerObj := model.Viewer{
		ID:        "11111",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		NotificationPreferences: model.ViewerNotificationPreferences{
			ID: "55555",
		},
		DiscussionID:     "22222",
		Discussion:       &model.Discussion{},
		LastViewed:       &lastViewed,
		LastViewedPostID: aws.String("33333"),
		LastViewedPost:   &model.Post{},
		Bookmarks:        &model.PostsConnection{},
		UserID:           "44444",
		User:             &model.User{},
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
				viewer: haveViewerObj,
			},
			want: map[string]*dynamodb.AttributeValue{
				"ViewerID": {
					S: aws.String(haveViewerObj.ID),
				},
				"CreatedAt": {
					S: aws.String(haveViewerObj.CreatedAt.Format(time.RFC3339Nano)),
				},
				"UpdatedAt": {
					S: aws.String(haveViewerObj.UpdatedAt.Format(time.RFC3339Nano)),
				},
				"DeletedAt": {
					NULL: aws.Bool(true),
				},
				"NotificationPreferences": {
					M: map[string]*dynamodb.AttributeValue{
						"ID": {
							S: aws.String(haveViewerObj.NotificationPreferences.ID),
						},
					},
				},
				"DiscussionID": {
					S: aws.String(haveViewerObj.DiscussionID),
				},
				"LastViewed": {
					S: aws.String(haveViewerObj.LastViewed.Format(time.RFC3339Nano)),
				},
				"LastViewedPostID": {
					S: haveViewerObj.LastViewedPostID,
				},
				"UserID": {
					S: aws.String(haveViewerObj.UserID),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshaled, err := datastoreObj.marshalMap(tt.args.viewer)
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
