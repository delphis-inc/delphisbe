package backend

import (
	"context"
	"fmt"
	"strings"

	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"go.uber.org/multierr"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *delphisBackend) GetDiscussionInviteByID(ctx context.Context, id string) (*model.DiscussionInvite, error) {
	return d.db.GetDiscussionInviteByID(ctx, id)
}

func (d *delphisBackend) GetDiscussionRequestAccessByID(ctx context.Context, id string) (*model.DiscussionAccessRequest, error) {
	return d.db.GetDiscussionRequestAccessByID(ctx, id)
}

func (d *delphisBackend) GetDiscussionInvitesByUserIDAndStatus(ctx context.Context, userID string, status model.InviteRequestStatus) ([]*model.DiscussionInvite, error) {
	iter := d.db.GetDiscussionInvitesByUserIDAndStatus(ctx, userID, status)
	return d.db.DiscussionInviteIterCollect(ctx, iter)
}

func (d *delphisBackend) GetSentDiscussionInvitesByUserID(ctx context.Context, userID string) ([]*model.DiscussionInvite, error) {
	iter := d.db.GetSentDiscussionInvitesByUserID(ctx, userID)
	return d.db.DiscussionInviteIterCollect(ctx, iter)
}

func (d *delphisBackend) GetDiscussionAccessRequestsByDiscussionID(ctx context.Context, discussionID string) ([]*model.DiscussionAccessRequest, error) {
	iter := d.db.GetDiscussionAccessRequestsByDiscussionID(ctx, discussionID)
	return d.db.AccessRequestIterCollect(ctx, iter)
}

func (d *delphisBackend) GetSentDiscussionAccessRequestsByUserID(ctx context.Context, userID string) ([]*model.DiscussionAccessRequest, error) {
	iter := d.db.GetSentDiscussionAccessRequestsByUserID(ctx, userID)
	return d.db.AccessRequestIterCollect(ctx, iter)
}

func (d *delphisBackend) GetInviteLinksByDiscussionID(ctx context.Context, discussionID string) (*model.DiscussionLinkAccess, error) {
	return d.db.GetInviteLinksByDiscussionID(ctx, discussionID)
}

func (d *delphisBackend) InviteUserToDiscussion(ctx context.Context, userID, discussionID, invitingParticipantID string) (*model.DiscussionInvite, error) {
	invite := model.DiscussionInvite{
		ID:                    util.UUIDv4(),
		UserID:                userID,
		DiscussionID:          discussionID,
		InvitingParticipantID: invitingParticipantID,
		Status:                model.InviteRequestStatusPending,
		InviteType:            model.InviteTypeInvite,
	}

	// TODO: Should block users from spamming invites?
	// Begin tx
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	// Put invite record
	inviteObj, err := d.db.PutDiscussionInviteRecord(ctx, tx, invite)
	if err != nil {
		logrus.WithError(err).Error("failed to put invite record")

		// Rollback on errors
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}
		return nil, err
	}

	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	return inviteObj, nil
}

func (d *delphisBackend) InviteTwitterUserToDiscussion(ctx context.Context, twitterHandle, discussionID, invitingParticipantID string) (*model.DiscussionInvite, error) {
	/* Get the authed user */
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	/* Obtain authed user profile */
	authedUserProfile, err := d.GetUserProfileByUserID(ctx, authedUser.UserID)
	if err != nil {
		return nil, err
	}

	/* Obtain authed user social info (this is needed for Twitter tokens) */
	authedSocialInfo, err := d.GetSocialInfosByUserProfileID(ctx, *&authedUserProfile.ID)
	if err != nil {
		return nil, err
	}

	/* Obtain infos needed for Twitter API client */
	consumerKey := d.config.Twitter.ConsumerKey
	consumerSecret := d.config.Twitter.ConsumerSecret
	accessToken := ""
	accessTokenSecret := ""
	for _, info := range authedSocialInfo {
		if strings.ToLower(info.Network) == "twitter" {
			accessToken = info.AccessToken
			accessTokenSecret = info.AccessTokenSecret
		}
	}

	if len(consumerKey) == 0 || len(consumerSecret) == 0 || len(accessToken) == 0 || len(accessTokenSecret) == 0 {
		return nil, fmt.Errorf("There is a problem retrieving authed user Twitter data")
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)

	/* Obtain invited user info from Twitter */
	httpClient := config.Client(oauth1.NoContext, token)
	twitterClient := twitter.NewClient(httpClient)
	userShowParams := &twitter.UserShowParams{ScreenName: twitterHandle}
	twitterUser, _, err := twitterClient.Users.Show(userShowParams)
	if err != nil {
		return nil, err
	}

	/* Check that users don't invite themselves */
	if twitterUser.ScreenName == authedUserProfile.TwitterHandle {
		return nil, fmt.Errorf("You cannot invite yourself")
	}

	/* Get invited user. If the user is not present in the system, we create it
	   with a dummy access token. Note, the system will not overwrite the tokens
	   with the dummy ones if valid tokens are already present */
	userObj, err := d.GetOrCreateUser(ctx, LoginWithTwitterInput{
		User:              twitterUser,
		AccessToken:       "",
		AccessTokenSecret: "",
	})
	if err != nil {
		logrus.WithError(err).Errorf("Got an error creating a user")
		return nil, err
	}

	/* Verify that an invite is not already present for such an user
	   NOTE: Should we check for already accepted invitations too? Maybe we can check if the user
	         is already a participant even before calling this function.*/
	userInvites, err := d.GetDiscussionInvitesByUserIDAndStatus(ctx, userObj.ID, model.InviteRequestStatusPending)
	if err != nil {
		return nil, err
	}
	if len(userInvites) > 0 {
		return nil, fmt.Errorf("This user already has a pending invitation")
	}

	/* Create a new invite for the user in the discussion */
	invite := model.DiscussionInvite{
		ID:                    util.UUIDv4(),
		UserID:                userObj.ID,
		DiscussionID:          discussionID,
		InvitingParticipantID: invitingParticipantID,
		Status:                model.InviteRequestStatusPending,
		InviteType:            model.InviteTypeInvite,
	}

	// Begin tx
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	// Put invite record
	inviteObj, err := d.db.PutDiscussionInviteRecord(ctx, tx, invite)
	if err != nil {
		logrus.WithError(err).Error("failed to put invite record")

		// Rollback on errors
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}
		return nil, err
	}

	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	/* TODO: (?) We may consider to notify users in some way external to the app, like email (if public) or twitter
	   dm (if they follow the authed user), in order to invite users to install the app. */

	return inviteObj, nil
}

func (d *delphisBackend) RequestAccessToDiscussion(ctx context.Context, userID, discussionID string) (*model.DiscussionAccessRequest, error) {
	request := model.DiscussionAccessRequest{
		ID:           util.UUIDv4(),
		UserID:       userID,
		DiscussionID: discussionID,
		Status:       model.InviteRequestStatusPending,
	}

	// TODO: Should block users from spamming requests?
	// Begin tx
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	// Put invite record
	requestObj, err := d.db.PutDiscussionAccessRequestRecord(ctx, tx, request)
	if err != nil {
		logrus.WithError(err).Error("failed to put request record")

		// Rollback on errors
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}
		return nil, err
	}

	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	return requestObj, nil
}

func (d *delphisBackend) RespondToInvitation(ctx context.Context, inviteID string, response model.InviteRequestStatus, discussionParticipantInput model.AddDiscussionParticipantInput) (*model.DiscussionInvite, error) {
	invite := model.DiscussionInvite{
		ID:     inviteID,
		Status: response,
	}

	// Begin tx
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	// Update invite record
	inviteObj, err := d.db.UpdateDiscussionInviteRecord(ctx, tx, invite)
	if err != nil {
		logrus.WithError(err).Error("failed to update invite record")

		// Rollback on errors
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}
		return nil, err
	}

	// If user has accepted the request, update discussion_user_access_table and create new participant
	if response == model.InviteRequestStatusAccepted {
		if _, err := d.db.UpsertDiscussionUserAccess(ctx, tx, inviteObj.DiscussionID, inviteObj.UserID); err != nil {
			logrus.WithError(err).Error("failed to update user access")

			// Rollback on errors
			if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
				logrus.WithError(txErr).Error("failed to rollback tx")
				return nil, multierr.Append(err, txErr)
			}
			return nil, err
		}

		if _, err := d.CreateParticipantForDiscussion(ctx, inviteObj.DiscussionID, inviteObj.UserID, discussionParticipantInput); err != nil {
			logrus.WithError(err).Error("failed to create participant for discussion")

			// Rollback on errors
			if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
				logrus.WithError(txErr).Error("failed to rollback tx")
				return nil, multierr.Append(err, txErr)
			}
			return nil, err
		}
	}

	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	return inviteObj, nil
}

func (d *delphisBackend) RespondToRequestAccess(ctx context.Context, requestID string, response model.InviteRequestStatus, invitingParticipantID string) (*model.DiscussionAccessRequest, error) {
	request := model.DiscussionAccessRequest{
		ID:     requestID,
		Status: response,
	}

	// Begin tx
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	// Update invite record
	requestObj, err := d.db.UpdateDiscussionAccessRequestRecord(ctx, tx, request)
	if err != nil {
		logrus.WithError(err).Error("failed to update request record")

		// Rollback on errors
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}
		return nil, err
	}

	// If user has accepted the request, update discussion_user_access_table to allow user to create participant when they join.
	if response == model.InviteRequestStatusAccepted {
		if _, err := d.db.UpsertDiscussionUserAccess(ctx, tx, requestObj.DiscussionID, requestObj.UserID); err != nil {
			logrus.WithError(err).Error("failed to update user access")

			// Rollback on errors
			if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
				logrus.WithError(txErr).Error("failed to rollback tx")
				return nil, multierr.Append(err, txErr)
			}
			return nil, err
		}

		// Send invite to user after they are accepted
		invite := model.DiscussionInvite{
			ID:                    util.UUIDv4(),
			UserID:                requestObj.UserID,
			DiscussionID:          requestObj.DiscussionID,
			InvitingParticipantID: invitingParticipantID,
			Status:                model.InviteRequestStatusPending,
			InviteType:            model.InviteTypeAccessRequestAccepted,
		}

		// Put invite record
		_, err := d.db.PutDiscussionInviteRecord(ctx, tx, invite)
		if err != nil {
			logrus.WithError(err).Error("failed to put invite record")
			// Rollback on errors
			if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
				logrus.WithError(txErr).Error("failed to rollback tx")
				return nil, multierr.Append(err, txErr)
			}
			return nil, err
		}
	}

	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	return requestObj, nil
}

func (d *delphisBackend) UpsertInviteLinksByDiscussionID(ctx context.Context, discussionID string) (*model.DiscussionLinkAccess, error) {
	input := model.DiscussionLinkAccess{
		DiscussionID:      discussionID,
		InviteLinkSlug:    util.UUIDv4(),
		VipInviteLinkSlug: util.UUIDv4(),
	}

	// Begin tx
	tx, err := d.db.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to begin tx")
		return nil, err
	}

	dla, err := d.db.UpsertInviteLinksByDiscussionID(ctx, tx, input)
	if err != nil {
		logrus.WithError(err).Error("failed to upsert discussion invite links")
		// Rollback on errors
		if txErr := d.db.RollbackTx(ctx, tx); txErr != nil {
			logrus.WithError(txErr).Error("failed to rollback tx")
			return nil, multierr.Append(err, txErr)
		}

		return nil, err
	}

	// Commit transaction
	if err := d.db.CommitTx(ctx, tx); err != nil {
		logrus.WithError(err).Error("failed to commit post tx")
		return nil, err
	}

	return dla, nil
}
