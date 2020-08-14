package test_utils

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/util"
)

const ProfileID = "profileID"
const DisplayName = "displayName"
const UserID = "userID"
const DiscussionID = "discussionID"
const LinkSlug = "slug"
const TwitterHandle = "twitterHandle"
const Token = "token"
const TokenSecret = "secret"
const Limit = 10
const ParticipantID = "participantID"
const ViewerID = "viewerID"
const PostID = "postID"
const PostContentID = "postContentID"
const ModeratorID = "modID"
const RequestID = "requestID"
const InvitingParticipantID = "invite_participating_id"
const GradientColor = model.GradientColorAzalea
const AnonymityType = model.AnonymityTypeStrong

var Now = time.Now()

func TestUser() model.User {
	return model.User{
		ID: UserID,
	}
}

func TestUserProfile() model.UserProfile {
	userID := UserID
	return model.UserProfile{
		ID:            ProfileID,
		DisplayName:   DisplayName,
		UserID:        &userID,
		TwitterHandle: TwitterHandle,
	}
}

func TestParticipant() model.Participant {
	discussionID := DiscussionID
	viewerID := ViewerID
	gradientColor := GradientColor
	userID := UserID

	return model.Participant{
		ID:            ParticipantID,
		ParticipantID: 0,
		DiscussionID:  &discussionID,
		ViewerID:      &viewerID,
		GradientColor: &gradientColor,
		UserID:        &userID,
		HasJoined:     false,
		IsAnonymous:   false,
	}
}

func TestDiscussion() model.Discussion {
	modID := ModeratorID

	return model.Discussion{
		ID:            DiscussionID,
		Title:         "title",
		AnonymityType: AnonymityType,
		ModeratorID:   &modID,
	}
}

func TestDiscussionShuffleTime() model.DiscussionShuffleTime {
	return model.DiscussionShuffleTime{
		DiscussionID: DiscussionID,
		ShuffleTime:  &Now,
	}
}

func TestDiscussionCreationSettings() model.DiscussionCreationSettings {
	return model.DiscussionCreationSettings{
		DiscussionJoinability: model.DiscussionJoinabilitySettingAllowTwitterFriends,
	}
}

func TestPost() model.Post {
	discussionID := DiscussionID
	participantID := ParticipantID
	postContentID := PostContentID

	return model.Post{
		ID:            PostID,
		PostType:      model.PostTypeStandard,
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContentID: &postContentID,
	}
}

func TestDiscussionArchive() model.DiscussionArchive {
	post := TestPost()

	postBytes, _ := json.Marshal([]*model.Post{&post})

	return model.DiscussionArchive{
		DiscussionID: DiscussionID,
		Archive:      postgres.Jsonb{postBytes},
	}
}

func TestPostContent() model.PostContent {
	return model.PostContent{
		ID:      PostContentID,
		Content: "hello world",
	}
}

func TestPostContentInput() model.PostContentInput {
	return model.PostContentInput{
		PostText: "hello world",
		PostType: model.PostTypeStandard,
	}
}

func TestModerator() model.Moderator {
	profileID := ProfileID

	return model.Moderator{
		ID:            ModeratorID,
		UserProfileID: &profileID,
		Discussion:    nil,
	}
}

func TestViewer() model.Viewer {
	discussionID := DiscussionID
	postID := PostID
	userID := UserID

	return model.Viewer{
		ID:               ViewerID,
		DiscussionID:     &discussionID,
		LastViewedPostID: &postID,
		UserID:           &userID,
	}
}

func TestSocialInfo() model.SocialInfo {
	return model.SocialInfo{
		Network:           util.SocialNetworkTwitter,
		AccessToken:       Token,
		AccessTokenSecret: TokenSecret,
		UserID:            UserID,
		UserProfileID:     ProfileID,
	}
}

func TestDiscussionInput() model.DiscussionInput {
	anonymityType := model.AnonymityTypeStrong
	title := "test title"
	publicAccess := true
	iconUrl := "http://test.com"

	return model.DiscussionInput{
		AnonymityType: &anonymityType,
		Title:         &title,
		PublicAccess:  &publicAccess,
		IconURL:       &iconUrl,
	}
}

func TestDiscussionsConnection() model.DiscussionsConnection {
	return model.DiscussionsConnection{
		IDs:   []string{DiscussionID},
		From:  0,
		To:    0,
		Edges: nil,
	}
}

func TestDiscussionAccessRequest(status model.InviteRequestStatus) model.DiscussionAccessRequest {
	return model.DiscussionAccessRequest{
		ID:           RequestID,
		UserID:       UserID,
		DiscussionID: DiscussionID,
		Status:       status,
	}
}

func TestDiscussionAccessLink() model.DiscussionAccessLink {
	return model.DiscussionAccessLink{
		DiscussionID: DiscussionID,
		LinkSlug:     LinkSlug,
	}
}

func TestAddDiscussionParticipantInput() model.AddDiscussionParticipantInput {
	gradientColor := GradientColor
	hasJoined := true

	return model.AddDiscussionParticipantInput{
		GradientColor: &gradientColor,
		HasJoined:     &hasJoined,
		IsAnonymous:   false,
	}
}

func TestUpdateParticipantInput() model.UpdateParticipantInput {
	gradientColor := GradientColor
	hasJoined := true
	isAnonymous := false

	return model.UpdateParticipantInput{
		GradientColor:   &gradientColor,
		HasJoined:       &hasJoined,
		IsAnonymous:     &isAnonymous,
		IsUnsetGradient: nil,
	}
}

func TestDiscussionUserAccess() model.DiscussionUserAccess {
	return model.DiscussionUserAccess{
		DiscussionID: DiscussionID,
		UserID:       UserID,
		State:        model.DiscussionUserAccessStateActive,
		NotifSetting: model.DiscussionUserNotificationSettingEverything,
	}
}

func TestPostsConnection(cursor string) model.PostsConnection {
	return model.PostsConnection{
		Edges: []*model.PostsEdge{
			{
				Cursor: cursor,
				Node:   &model.Post{},
			},
		},
		PageInfo: model.PageInfo{},
	}
}

func TestDelphisAuthedUser() auth.DelphisAuthedUser {
	user := TestUser()
	user.ID = "authoredUserID"
	return auth.DelphisAuthedUser{
		UserID: user.ID,
		User:   &user,
	}
}
