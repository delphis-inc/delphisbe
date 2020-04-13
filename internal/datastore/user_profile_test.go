package datastore

import (
	"testing"
)

func Test_MarshalUserProfile(t *testing.T) {
	// type args struct {
	// 	userProfile model.UserProfile
	// }

	// haveUserProfileObj := model.UserProfile{
	// 	ID:                     "11111",
	// 	DisplayName:            "Delphis Hello",
	// 	UserID:                 "22222",
	// 	TwitterHandle:          "delphishq",
	// 	ModeratedDiscussionIDs: []string{"33333"},
	// 	ModeratedDiscussions:   []model.Discussion{model.Discussion{}},
	// 	TwitterInfo: model.SocialInfo{
	// 		AccessToken:       "44444",
	// 		AccessTokenSecret: "55555",
	// 		UserID:            "55555",
	// 		ProfileImageURL:   "https://a.b/c.png",
	// 		ScreenName:        "delphishq",
	// 		IsVerified:        true,
	// 	},
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
	// 			userProfile: haveUserProfileObj,
	// 		},
	// 		want: map[string]*dynamodb.AttributeValue{
	// 			"ID": {
	// 				S: aws.String(haveUserProfileObj.ID),
	// 			},
	// 			"DisplayName": {
	// 				S: aws.String(haveUserProfileObj.DisplayName),
	// 			},
	// 			"UserID": {
	// 				S: aws.String(haveUserProfileObj.UserID),
	// 			},
	// 			"TwitterHandle": {
	// 				S: aws.String(haveUserProfileObj.TwitterHandle),
	// 			},
	// 			"ModeratedDiscussionIDs": {
	// 				SS: []*string{aws.String(haveUserProfileObj.ModeratedDiscussionIDs[0])},
	// 			},
	// 			"TwitterInfo": {
	// 				M: map[string]*dynamodb.AttributeValue{
	// 					"AccessToken": {
	// 						S: aws.String(haveUserProfileObj.TwitterInfo.AccessToken),
	// 					},
	// 					"AccessTokenSecret": {
	// 						S: aws.String(haveUserProfileObj.TwitterInfo.AccessTokenSecret),
	// 					},
	// 					"SocialUserID": {
	// 						S: aws.String(haveUserProfileObj.TwitterInfo.UserID),
	// 					},
	// 					"ProfileImageURL": {
	// 						S: aws.String(haveUserProfileObj.TwitterInfo.ProfileImageURL),
	// 					},
	// 					"ScreenName": {
	// 						S: aws.String(haveUserProfileObj.TwitterInfo.ScreenName),
	// 					},
	// 					"IsVerified": {
	// 						BOOL: aws.Bool(haveUserProfileObj.TwitterInfo.IsVerified),
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		marshaled, err := datastoreObj.marshalMap(tt.args.userProfile)
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
