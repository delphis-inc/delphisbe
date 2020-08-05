package test_utils

import (
	"time"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/util"
)

const ProfileID = "profileID"
const DisplayName = "displayName"
const UserID = "userID"
const InviterID = "InviterID"
const DiscussionID = "discussionID"
const LinkSlug = "slug"
const TwitterHandle = "twitterHandle"
const Token = "token"
const TokenSecret = "secret"
const Limit = 10
const ContentID = "contentID"
const ContentName = "contentName"
const ContentType = "contentType"
const ContentLink = "http://content.link"
const IdleMinutes = 120
const Email = "test@email.com"
const ParticipantID = "participantID"
const ViewerID = "viewerID"
const FlairID = "flairID"
const FlairTemplateID = "templateID"
const PostID = "postID"
const PostContentID = "postContentID"
const ModeratorID = "modID"
const Tag = "tag"
const Source = "source"
const ImageURL = "http://image.url"
const InviteID = "inviteID"
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
	flairID := FlairID
	gradientColor := GradientColor
	userID := UserID
	inviterID := InviterID

	return model.Participant{
		ID:            ParticipantID,
		ParticipantID: 0,
		DiscussionID:  &discussionID,
		ViewerID:      &viewerID,
		FlairID:       &flairID,
		GradientColor: &gradientColor,
		UserID:        &userID,
		HasJoined:     false,
		IsAnonymous:   false,
		InviterID:     &inviterID,
	}
}

func TestDiscussion() model.Discussion {
	modID := ModeratorID

	return model.Discussion{
		ID:            DiscussionID,
		Title:         "title",
		AnonymityType: AnonymityType,
		ModeratorID:   &modID,
		AutoPost:      false,
		IdleMinutes:   IdleMinutes,
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

func TestContentQueueRecord() model.ContentQueueRecord {
	return model.ContentQueueRecord{
		DiscussionID:      DiscussionID,
		ImportedContentID: ContentID,
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

func TestFlair() model.Flair {
	return model.Flair{
		ID:         FlairID,
		TemplateID: FlairTemplateID,
		UserID:     UserID,
	}
}

func TestFlairTemplate() model.FlairTemplate {
	displayName := DisplayName
	imageURL := ImageURL
	return model.FlairTemplate{
		ID:          FlairTemplateID,
		DisplayName: &displayName,
		ImageURL:    &imageURL,
		Source:      Source,
	}
}

func TestDiscussionAutoPost() model.DiscussionAutoPost {
	return model.DiscussionAutoPost{
		ID:          DiscussionID,
		IdleMinutes: 120,
	}
}

func TestImportedContent() model.ImportedContent {
	return model.ImportedContent{
		ID:          ContentID,
		ContentName: ContentName,
		ContentType: ContentType,
		Link:        ContentLink,
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
	autoPost := true
	idleMinutes := IdleMinutes
	publicAccess := true
	iconUrl := "http://test.com"

	return model.DiscussionInput{
		AnonymityType: &anonymityType,
		Title:         &title,
		AutoPost:      &autoPost,
		IdleMinutes:   &idleMinutes,
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

func TestDiscussionTag() model.Tag {
	return model.Tag{
		ID:  DiscussionID,
		Tag: Tag,
	}
}

func TestContentTag() model.Tag {
	return model.Tag{
		ID:  ContentID,
		Tag: Tag,
	}
}

func TestDiscussionFlairTemplateAccess() model.DiscussionFlairTemplateAccess {
	return model.DiscussionFlairTemplateAccess{
		DiscussionID:    DiscussionID,
		FlairTemplateID: FlairTemplateID,
	}
}

func TestImportedContentInput() model.ImportedContentInput {
	return model.ImportedContentInput{

		ContentName: ContentName,
		ContentType: ContentType,
		Link:        ContentLink,
		Overview:    "",
		Source:      "",
	}
}

func TestDiscussionInvite(status model.InviteRequestStatus) model.DiscussionInvite {
	return model.DiscussionInvite{
		ID:                    InviteID,
		UserID:                UserID,
		DiscussionID:          DiscussionID,
		InvitingParticipantID: InvitingParticipantID,
		Status:                status,
		InviteType:            "",
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
	flairID := FlairID
	hasJoined := true

	return model.AddDiscussionParticipantInput{
		GradientColor: &gradientColor,
		FlairID:       &flairID,
		HasJoined:     &hasJoined,
		IsAnonymous:   false,
	}
}

func TestUpdateParticipantInput() model.UpdateParticipantInput {
	gradientColor := GradientColor
	flairID := FlairID
	hasJoined := true
	isAnonymous := false

	return model.UpdateParticipantInput{
		GradientColor:   &gradientColor,
		FlairID:         &flairID,
		HasJoined:       &hasJoined,
		IsAnonymous:     &isAnonymous,
		IsUnsetGradient: nil,
		IsUnsetFlairID:  nil,
	}
}

func TestDiscussionUserAccess() model.DiscussionUserAccess {
	return model.DiscussionUserAccess{
		DiscussionID: DiscussionID,
		UserID:       UserID,
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
