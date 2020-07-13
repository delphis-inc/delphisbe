package datastore

import (
	"context"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) GetParticipantByID(ctx context.Context, id string) (*model.Participant, error) {
	logrus.Debugf("GetParticipantByID::SQL Query")
	participants, err := d.GetParticipantsByIDs(ctx, []string{id})
	if err != nil {
		return nil, err
	}

	return participants[id], nil
}

func (d *delphisDB) GetParticipantsByIDs(ctx context.Context, ids []string) (map[string]*model.Participant, error) {
	logrus.Debugf("GetParticipantsByIDs::SQL Query")
	participants := []model.Participant{}
	if err := d.sql.Where(ids).Find(&participants).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetParticipantByID::Failed to query participant by ID")
		return nil, err
	}

	retVal := map[string]*model.Participant{}
	for _, p := range participants {
		tempVal := p
		retVal[p.ID] = &tempVal
	}

	return retVal, nil
}

func (d *delphisDB) GetParticipantsByDiscussionID(ctx context.Context, id string) ([]model.Participant, error) {
	logrus.Debugf("GetParticipantsByDiscussionID::SQL Query")
	participants := []model.Participant{}
	if err := d.sql.Where(&model.Participant{DiscussionID: &id}).Order("participant_id desc").Find(&participants).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetParticipantsByDiscussionID::Failed to get participants by discussion ID")
		return nil, err
	}
	return participants, nil
}

func (d *delphisDB) GetParticipantsByDiscussionIDUserID(ctx context.Context, discussionID string, userID string) ([]model.Participant, error) {
	logrus.Debugf("GetParticipantByDiscussionIDUserID::SQL Query")
	participants := []model.Participant{}
	if err := d.sql.Where(&model.Participant{DiscussionID: &discussionID, UserID: &userID}).Order("participant_id desc").Limit(2).Find(&participants).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetParticipantByDiscussionIDUserID::Failed to get participant by discussion ID and user ID")
		return nil, err
	}
	return participants, nil
}

func (d *delphisDB) GetModeratorParticipantsByDiscussionID(ctx context.Context, discussionID string) ([]model.Participant, error) {
	logrus.Debugf("GetModeratorParticipantsByDiscussionID::SQL Query")
	participants := []model.Participant{}
	joinUserProfiles := "JOIN user_profiles ON user_profiles.user_id = participants.user_id"
	joinModerators := "JOIN moderators ON user_profiles.id = moderators.user_profile_id"
	joinDiscussions := "JOIN discussions ON participants.discussion_id = discussions.id"
	if err := d.sql.Joins(joinUserProfiles).Joins(joinModerators).Joins(joinDiscussions).Where(`("discussions"."moderator_id" = "moderators"."id") AND "discussions"."id" = ?`, &discussionID).Order("participant_id desc").Limit(2).Find(&participants).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetModeratorParticipantsByDiscussionID::Failed to get moderator participants by discussion ID")
		return nil, err
	}
	return participants, nil
}

func (d *delphisDB) UpsertParticipant(ctx context.Context, participant model.Participant) (*model.Participant, error) {
	logrus.Debug("UpsertParticipant::SQL Create")
	found := model.Participant{}
	if err := d.sql.First(&found, model.Participant{ID: participant.ID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := d.sql.Create(&participant).First(&found, model.Participant{ID: participant.ID}).Error; err != nil {
				logrus.WithError(err).Errorf("UpsertParticipant::Failed to put Participant")
				return nil, err
			}
		} else {
			logrus.WithError(err).Errorf("UpsertParticipant::Failed checking for Participant object")
			return nil, err
		}
	} else {
		if err := d.sql.Model(&participant).Updates(model.Participant{
			FlairID:       participant.FlairID,
			IsAnonymous:   participant.IsAnonymous,
			UpdatedAt:     time.Now(),
			GradientColor: participant.GradientColor,
			HasJoined:     participant.HasJoined,
			IsBanned:      participant.IsBanned,
		}).First(&found).Error; err != nil {
			logrus.WithError(err).Errorf("UpsertParticipant::Failed updating Participant object")
			return nil, err
		}
	}
	return &found, nil
}

func (d *delphisDB) AssignFlair(ctx context.Context, participant model.Participant, flairID *string) (*model.Participant, error) {
	logrus.Debug("AssignFlair::SQL Update")
	if err := d.sql.Model(&participant).UpdateColumn("FlairID", flairID).Error; err != nil {
		logrus.WithError(err).Errorf("AssignFlair::Failed to update")
		return &participant, err
	}
	return &participant, nil
}

func (d *delphisDB) GetTotalParticipantCountByDiscussionID(ctx context.Context, discussionID string) int {
	count := 0
	d.sql.Model(&model.Participant{}).Where(&model.Participant{DiscussionID: &discussionID}).Count(&count)
	return count
}
