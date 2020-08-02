package backend

import (
	"context"

	"github.com/delphis-inc/delphisbe/internal/util"
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

func (d *delphisBackend) GetDiscussionAccessRequestByDiscussionIDUserID(ctx context.Context, discussionID string, userID string) (*model.DiscussionAccessRequest, error) {
	return d.db.GetDiscussionAccessRequestByDiscussionIDUserID(ctx, discussionID, userID)
}

func (d *delphisBackend) GetSentDiscussionAccessRequestsByUserID(ctx context.Context, userID string) ([]*model.DiscussionAccessRequest, error) {
	iter := d.db.GetSentDiscussionAccessRequestsByUserID(ctx, userID)
	return d.db.AccessRequestIterCollect(ctx, iter)
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
