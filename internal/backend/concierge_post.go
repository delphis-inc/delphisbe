package backend

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/nedrocks/delphisbe/internal/util"

	"github.com/nedrocks/delphisbe/graph/model"
)

func (d *delphisBackend) HandleConciergeMutation(ctx context.Context, userID string, discussionID string, mutationID string, selectedOptions []string) (bool, error) {
	switch mutationID {
	case string(model.MutationUpdateFlairAccessToDiscussion):
		logrus.Debugf("Update flair access: %v\n", selectedOptions)
		return d.updateDiscussionFlairAccess(ctx, discussionID, selectedOptions)
	case string(model.MutationUpdateDiscussionNameAndEmoji):
		logrus.Debugf("Update discussion name and emoji: %v\n", selectedOptions)
	case string(model.MutationUpdateViewerAccessibility):
		logrus.Debugf("Update viewer accessibility: %v\n", selectedOptions)
	case string(model.MutationUpdateInvitationApproval):
		logrus.Debugf("Update invitation approval: %v\n", selectedOptions)
	}
	return false, nil
}

func (d *delphisBackend) createInviteLinkConciergePost(ctx context.Context, discussionID string, participantID string) (*model.Post, error) {
	// Get invite links by discussionID
	// TODO: Finish after merging inviteLinkPR

	appAction := string(model.AppActionCopyToClipboard)
	content := model.ConciergeContent{
		AppActionID: &appAction,
		Options: []*model.ConciergeOption{
			{
				Text:  "Copy VIP Link (auto-join)",
				Value: "https://delphis.com/viplink/testme", // update
			},
			{
				Text:  "Copy public Link (approval req'd)",
				Value: "https://delphis.com/publiclink/testme", // update
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
	}

	return &post, nil
}

func (d *delphisBackend) createInviteSettingConciergePost(ctx context.Context, discussionID, participantID string, discussionObj *model.Discussion) (*model.Post, error) {
	mutationID := string(model.MutationUpdateInvitationApproval)
	content := model.ConciergeContent{
		MutationID: &mutationID,
		Options: []*model.ConciergeOption{
			{
				Text:     "Manually approve all invites",
				Value:    "false",
				Selected: false, // update to discussionObj.AutoJoin
			},
			{
				Text:     "Invites get auto-joined",
				Value:    "true",
				Selected: true, // update to discussionObj.AutoJoin
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
	}

	return &post, nil
}

func (d *delphisBackend) createViewerAccessConciergePost(ctx context.Context, discussionID, participantID string, discussionObj *model.Discussion) (*model.Post, error) {
	mutationID := string(model.MutationUpdateViewerAccessibility)
	content := model.ConciergeContent{
		MutationID: &mutationID,
		Options: []*model.ConciergeOption{
			{
				Text:     "Make chat publicly viewable",
				Value:    "true",
				Selected: false, // update to discussionObj.PublicToViewer
			},
			{
				Text:     "Chat is private to participants",
				Value:    "false",
				Selected: true, // update to discussionObj.PublicToViewer
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
	}

	return &post, nil
}

func (d *delphisBackend) createRenameChatAndEmojiConciergePost(ctx context.Context, discussionID, participantID string) (*model.Post, error) {
	mutationID := string(model.MutationUpdateDiscussionNameAndEmoji)
	content := model.ConciergeContent{
		MutationID: &mutationID,
		Options: []*model.ConciergeOption{
			{
				Text:  "Rename chat + pick emoji",
				Value: "return chat name and emoji comma-separated", // how do we want to handle this on the client?
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
	}

	return &post, nil
}

func (d *delphisBackend) updateDiscussionFlairAccess(ctx context.Context, discussionID string, flairTemplates []string) (bool, error) {
	// TODO: need flair access PR to be merged

	return true, nil
}
