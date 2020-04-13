package model

import (
	"testing"
)

func Test_UserDynamoUnmarshal(t *testing.T) {
	// type args struct {
	// 	item map[string]*dynamodb.AttributeValue
	// }

	// tests := []struct {
	// 	name string
	// 	args args
	// 	want *User
	// }{
	// 	{
	// 		name: "single_values",
	// 		args: args{
	// 			map[string]*dynamodb.AttributeValue{
	// 				"DiscussionParticipants": {
	// 					SS: []*string{aws.String("5b4baec1-300d-420e-9dc8-a1eb0ee83d41.0")},
	// 				},
	// 				"DiscussionViewers": {
	// 					SS: []*string{aws.String("5b4baec1-300d-420e-9dc8-a1eb0ee83d41.7505e149-aa0d-48e8-a6ff-07f22f4a1709")},
	// 				},
	// 				"ID": {
	// 					S: aws.String("1a4deec5-ae6f-40d3-94db-15d5ba84166a"),
	// 				},
	// 				"createdAt": {
	// 					S: aws.String("2020-03-06T16:50:48.288033-08:00"),
	// 				},
	// 				"deletedAt": {
	// 					NULL: aws.Bool(true),
	// 				},
	// 				"participantIDs": {
	// 					NULL: aws.Bool(true),
	// 				},
	// 				"updatedAt": {
	// 					S: aws.String("2020-03-06T16:50:48.288033-08:00"),
	// 				},
	// 				"viewerIDs": {
	// 					NULL: aws.Bool(true),
	// 				},
	// 			},
	// 		},
	// 		want: &User{
	// 			DiscussionParticipants: DiscussionParticipantKeys{
	// 				Keys: []DiscussionParticipantKey{
	// 					{
	// 						DiscussionID:  "5b4baec1-300d-420e-9dc8-a1eb0ee83d41",
	// 						ParticipantID: 0,
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		userObj := User{}
	// 		err := dynamodbattribute.UnmarshalMap(tt.args.item, &userObj)
	// 		if err != nil {
	// 			t.Errorf("Caught an error: %+v", err)
	// 			return
	// 		}
	// 		if !reflect.DeepEqual(userObj.DiscussionParticipants, tt.want.DiscussionParticipants) {
	// 			t.Errorf("These objects did not match. Got: %+v; Want: %+v", userObj.DiscussionParticipants, tt.want.DiscussionParticipants)
	// 		}
	// 	})
	// }
}
