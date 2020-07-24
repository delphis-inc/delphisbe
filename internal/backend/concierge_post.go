package backend

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/delphis-inc/delphisbe/internal/util"

	"github.com/delphis-inc/delphisbe/graph/model"
)

func (d *delphisBackend) GetConciergeParticipantID(ctx context.Context, discussionID string) (string, error) {
	// Get concierge participant for posts
	participants, err := d.GetParticipantsByDiscussionIDUserID(ctx, discussionID, model.ConciergeUser)
	if err != nil {
		logrus.WithError(err).Error("failed to get concierge participant")
		return "", err
	}

	if participants.NonAnon == nil {
		return "", fmt.Errorf("no non-anonymous participant for the concierge")
	}

	return participants.NonAnon.ID, nil
}

func (d *delphisBackend) HandleConciergeMutation(ctx context.Context, userID string, discussionID string, mutationID string, selectedOptions []string) (*model.Post, error) {
	conciergeParticipant, err := d.GetConciergeParticipantID(ctx, discussionID)
	if err != nil {
		logrus.WithError(err).Error("failed to get concierge participant")
		return nil, err
	}

	switch mutationID {
	case string(model.MutationUpdateFlairAccessToDiscussion):
		logrus.Debugf("Update flair access: %v\n", selectedOptions)
		if err := d.updateDiscussionFlairAccess(ctx, discussionID, selectedOptions); err != nil {
			logrus.WithError(err).Error("failed to update discussion flair access")
			return nil, err
		}

		return d.createFlairAccessConciergePost(ctx, userID, discussionID, conciergeParticipant)
	case string(model.MutationUpdateDiscussionNameAndEmoji):
		logrus.Debugf("Update discussion name and emoji: %v\n", selectedOptions)

		return d.createRenameChatAndEmojiConciergePost(ctx, discussionID, conciergeParticipant)
	case string(model.MutationUpdateViewerAccessibility):
		logrus.Debugf("Update viewer accessibility: %v\n", selectedOptions)

		return d.createViewerAccessConciergePost(ctx, discussionID, conciergeParticipant)
	case string(model.MutationUpdateInvitationApproval):
		logrus.Debugf("Update invitation approval: %v\n", selectedOptions)

		return d.createInviteSettingConciergePost(ctx, discussionID, conciergeParticipant)
	}

	logrus.Error("mutationID did not match any existing mutations")
	return nil, fmt.Errorf("mutationID did not match any existing mutations")
}

func (d *delphisBackend) createInviteLinkConciergePost(ctx context.Context, discussionID string, participantID string) (*model.Post, error) {
	// Get invite links by discussionID

	links, err := d.GetInviteLinksByDiscussionID(ctx, discussionID)
	if err != nil || links == nil {
		return nil, err
	}

	appAction := string(model.AppActionCopyToClipboard)
	content := model.ConciergeContent{
		AppActionID: &appAction,
		Options: []*model.ConciergeOption{
			{
				Text:  "Copy VIP Link (auto-join)",
				Value: fmt.Sprintf("https://m.chatham.ai/d/%s/%s", discussionID, links.VipInviteLinkSlug),
			},
			{
				Text:  "Copy public Link (approval req'd)",
				Value: fmt.Sprintf("https://m.chatham.ai/d/%s", discussionID),
			},
		},
	}

	post := model.Post{
		ID:            util.UUIDv4(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			Content: "Here's your links for the chat. By default, you have to approve people who try to join, but you can share the VIP link too.",
		},
		ConciergeContent: &content,
		PostType:         model.PostTypeConcierge,
	}

	return &post, nil
}

func (d *delphisBackend) createFlairAccessConciergePost(ctx context.Context, userID, discussionID, participantID string) (*model.Post, error) {
	// Get user's flair
	flairs, err := d.GetFlairsByUserID(ctx, userID)
	if err != nil {
		logrus.WithError(err).Error("failed to get flairs by userID")
		return nil, err
	}

	var options []*model.ConciergeOption
	for _, flair := range flairs {
		flairObj, err := d.GetFlairTemplateByID(ctx, flair.TemplateID)
		if err != nil {
			logrus.WithError(err).Error("failed to get flair template by ID")
			return nil, err
		}

		tempOption := model.ConciergeOption{
			Text:  *flairObj.DisplayName,
			Value: flair.TemplateID,
		}
		options = append(options, &tempOption)
	}

	mutationID := string(model.MutationUpdateFlairAccessToDiscussion)
	content := model.ConciergeContent{
		MutationID: &mutationID,
		Options:    options,
	}

	post := model.Post{
		ID:            util.UUIDv4(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			Content: "You have some nice flair! If you like, you can set it up so people from the same groups can join automatically." +
				" Otherwise, they'll need the VIP link or you'll have to approve them before they can join the chat.",
		},
		ConciergeContent: &content,
		PostType:         model.PostTypeConcierge,
	}

	return &post, nil
}

func (d *delphisBackend) createInviteSettingConciergePost(ctx context.Context, discussionID, participantID string) (*model.Post, error) {
	// Get discussionObj to show which settings are selected
	_, err := d.GetDiscussionByID(ctx, discussionID)
	if err != nil {
		logrus.WithError(err).Error("failed to get discussion by ID")
		return nil, err
	}

	mutationID := string(model.MutationUpdateInvitationApproval)
	content := model.ConciergeContent{
		MutationID: &mutationID,
		Options: []*model.ConciergeOption{
			{
				Text:     "Manually approve all invites",
				Value:    "false",
				Selected: false, // TODO: update to discussionObj.AutoJoin
			},
			{
				Text:     "Invites get auto-joined",
				Value:    "true",
				Selected: true, // TODO: update to discussionObj.AutoJoin
			},
		},
	}

	post := model.Post{
		ID:            util.UUIDv4(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			Content: "Also, what should we do it someone in the chat wants to invite a friend?",
		},
		ConciergeContent: &content,
		PostType:         model.PostTypeConcierge,
	}

	return &post, nil
}

func (d *delphisBackend) createViewerAccessConciergePost(ctx context.Context, discussionID, participantID string) (*model.Post, error) {
	// Get discussionObj to show which settings are selected
	_, err := d.GetDiscussionByID(ctx, discussionID)
	if err != nil {
		logrus.WithError(err).Error("failed to get discussion by ID")
		return nil, err
	}

	mutationID := string(model.MutationUpdateViewerAccessibility)
	content := model.ConciergeContent{
		MutationID: &mutationID,
		Options: []*model.ConciergeOption{
			{
				Text:     "Make chat publicly viewable",
				Value:    "true",
				Selected: false, // TODO: update to discussionObj.PublicToViewer
			},
			{
				Text:     "Chat is private to participants",
				Value:    "false",
				Selected: true, // TODO: update to discussionObj.PublicToViewer
			},
		},
	}

	post := model.Post{
		ID:            util.UUIDv4(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			Content: "Do you want to allow people who aren't in the chat to be able to watch?",
		},
		ConciergeContent: &content,
		PostType:         model.PostTypeConcierge,
	}

	return &post, nil
}

func (d *delphisBackend) createRenameChatAndEmojiConciergePost(ctx context.Context, discussionID, participantID string) (*model.Post, error) {
	actionID := string(model.AppActionRenameChat)
	content := model.ConciergeContent{
		AppActionID: &actionID,
		Options: []*model.ConciergeOption{
			{
				Text:  "Rename chat + pick emoji",
				Value: "rename",
			},
		},
	}

	post := model.Post{
		ID:            util.UUIDv4(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DiscussionID:  &discussionID,
		ParticipantID: &participantID,
		PostContent: &model.PostContent{
			Content: "Would you like to rename the chat or pick a different emoji?", // Need copy
		},
		ConciergeContent: &content,
		PostType:         model.PostTypeConcierge,
	}

	return &post, nil
}

func (d *delphisBackend) updateDiscussionFlairAccess(ctx context.Context, discussionID string, flairTemplates []string) error {
	// TODO: need flair access PR to be merged

	return nil
}
