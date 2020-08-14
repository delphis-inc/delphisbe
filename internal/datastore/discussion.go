package datastore

import (
	"context"
	"database/sql"
	"io"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) GetDiscussionByID(ctx context.Context, id string) (*model.Discussion, error) {
	logrus.Debug("GetDiscussionByID::SQL Query")
	discussions, err := d.GetDiscussionsByIDs(ctx, []string{id})
	if err != nil {
		return nil, err
	}

	return discussions[id], nil
}

func (d *delphisDB) IncrementDiscussionShuffleCount(ctx context.Context, tx *sql.Tx, id string) (*int, error) {
	logrus.Debug("IncrementDiscussionShuffleCount::SQL Update")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("IncrementDiscussionShuffleCount::failed to initialize statements")
		return nil, err
	}

	discussion := model.Discussion{}
	if err := tx.StmtContext(ctx, d.prepStmts.incrDiscussionShuffleCount).QueryRowContext(
		ctx,
		id,
	).Scan(
		&discussion.ShuffleCount,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Errorf("failed to increment Shuffle ID for discussion with ID %s", id)
		return nil, err
	}

	return &discussion.ShuffleCount, nil
}

func (d *delphisDB) GetDiscussionByLinkSlug(ctx context.Context, slug string) (*model.Discussion, error) {
	logrus.Debug("GetDiscussionByLinkSlug::SQL Update")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetDiscussionByLinkSlug::failed to initialize statements")
		return nil, err
	}

	discussion := model.Discussion{}
	titleHistory := make([]byte, 0)
	descriptionHistory := make([]byte, 0)

	if err := d.prepStmts.getDiscussionByLinkSlugStmt.QueryRowContext(
		ctx,
		slug,
	).Scan(
		&discussion.ID,
		&discussion.CreatedAt,
		&discussion.UpdatedAt,
		&discussion.DeletedAt,
		&discussion.Title,
		&discussion.AnonymityType,
		&discussion.ModeratorID,
		&discussion.IconURL,
		&discussion.Description,
		&titleHistory,
		&descriptionHistory,
		&discussion.DiscussionJoinability,
		&discussion.LastPostID,
		&discussion.LastPostCreatedAt,
		&discussion.ShuffleCount,
		&discussion.LockStatus,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logrus.WithError(err).Error("failed to execute getDiscussionByLinkSlugStmt")
		return nil, err
	}

	discussion.TitleHistory.RawMessage = titleHistory
	discussion.DescriptionHistory.RawMessage = descriptionHistory

	return &discussion, nil
}

func (d *delphisDB) GetDiscussionsByIDs(ctx context.Context, ids []string) (map[string]*model.Discussion, error) {
	logrus.Debug("GetDiscussionsByIDs::SQL Query")
	discussions := []model.Discussion{}
	if err := d.sql.Where(ids).Preload("Moderator").Find(&discussions).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// This is a not found situation with multiple ids. Not sure what to do here..?
		} else {
			logrus.WithError(err).Errorf("GetDiscussionsByIDs::Failed to get discussions by IDs")
			return nil, err
		}
	}
	retVal := map[string]*model.Discussion{}
	for _, id := range ids {
		retVal[id] = nil
	}
	for _, disc := range discussions {
		tempDisc := disc
		retVal[disc.ID] = &tempDisc
	}
	return retVal, nil
}

func (d *delphisDB) GetDiscussionByModeratorID(ctx context.Context, moderatorID string) (*model.Discussion, error) {
	logrus.Debugf("GetDiscussionByModeratorID::SQL Query")
	discussion := model.Discussion{}
	if err := d.sql.First(&discussion, model.Discussion{ModeratorID: &moderatorID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetDiscussionByModeratorID::Failed getting discussion by moderator ID")
		return nil, err
	}

	return &discussion, nil
}

func (d *delphisDB) ListDiscussions(ctx context.Context) (*model.DiscussionsConnection, error) {
	//TODO: this should take in paging params and return based on those.
	logrus.Debugf("ListDiscussions::SQL Query")

	discussions := []model.Discussion{}
	if err := d.sql.Preload("Moderator").Find(&discussions).Error; err != nil {
		logrus.WithError(err).Errorf("ListDiscussions::Failed to list discussions")
		return nil, err
	}

	ids := make([]string, 0)
	edges := make([]*model.DiscussionsEdge, 0)
	for i := range discussions {
		discussionObj := &discussions[i]
		edges = append(edges, &model.DiscussionsEdge{
			Node: discussionObj,
		})
		ids = append(ids, discussionObj.ID)
	}

	return &model.DiscussionsConnection{
		Edges: edges,
		IDs:   ids,
	}, nil
}

func (d *delphisDB) ListDiscussionsByUserID(ctx context.Context, userID string, state model.DiscussionUserAccessState) (*model.DiscussionsConnection, error) {
	//TODO: this should take in paging params and return based on those.
	logrus.Debugf("ListDiscussions::SQL Query")

	logrus.Infof("State: %+v\n", state)

	iter := d.GetDiscussionsByUserAccess(ctx, userID, state)
	discArr, err := d.DiscussionIterCollect(ctx, iter)
	if err != nil {
		logrus.WithError(err).Error("failed to collect discussions from iter")
		return nil, err
	}

	ids := make([]string, 0)
	edges := make([]*model.DiscussionsEdge, 0)
	for i := range discArr {
		discussionObj := discArr[i]
		edges = append(edges, &model.DiscussionsEdge{
			Node: discussionObj,
		})
		ids = append(ids, discussionObj.ID)
	}

	return &model.DiscussionsConnection{
		Edges: edges,
		IDs:   ids,
	}, nil
}

func (d *delphisDB) UpsertDiscussion(ctx context.Context, discussion model.Discussion) (*model.Discussion, error) {
	logrus.Debug("UpsertDiscussion::SQL Create")
	found := model.Discussion{}
	if err := d.sql.First(&found, model.Discussion{ID: discussion.ID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if err := d.sql.Create(&discussion).First(&found, model.Discussion{ID: discussion.ID}).Error; err != nil {
				logrus.WithError(err).Errorf("UpsertDiscussion::Failed to create new object")
				return nil, err
			}
		} else {
			logrus.WithError(err).Errorf("UpsertDiscussion::Failed checking for Discussion object")
			return nil, err
		}
	} else {
		if err := d.sql.Preload("Moderator").Model(&discussion).Updates(map[string]interface{}{
			"Title":                 discussion.Title,
			"Description":           discussion.Description,
			"AnonymityType":         discussion.AnonymityType,
			"IconURL":               discussion.IconURL,
			"TitleHistory":          discussion.TitleHistory,
			"DescriptionHistory":    discussion.DescriptionHistory,
			"DiscussionJoinability": discussion.DiscussionJoinability,
			"LastPostID":            discussion.LastPostID,
			"LastPostCreatedAt":     discussion.LastPostCreatedAt,
			"LockStatus":            discussion.LockStatus,
		}).First(&found).Error; err != nil {
			logrus.WithError(err).Errorf("UpsertDiscussion::Failed updating disucssion object")
			return nil, err
		}
	}
	return &found, nil
}

func (d *delphisDB) DiscussionIterCollect(ctx context.Context, iter DiscussionIter) ([]*model.Discussion, error) {
	var discussions []*model.Discussion
	disc := model.Discussion{}

	defer iter.Close()

	for iter.Next(&disc) {
		tempDisc := disc

		discussions = append(discussions, &tempDisc)
	}

	if err := iter.Close(); err != nil && err != io.EOF {
		logrus.WithError(err).Error("failed to close iter")
		return nil, err
	}

	return discussions, nil
}

type discussionIter struct {
	err  error
	ctx  context.Context
	rows *sql.Rows
}

func (iter *discussionIter) Next(discussion *model.Discussion) bool {
	if iter.err != nil {
		logrus.WithError(iter.err).Error("iterator error")
		return false
	}

	if iter.err = iter.ctx.Err(); iter.err != nil {
		logrus.WithError(iter.err).Error("iterator context error")
		return false
	}

	if !iter.rows.Next() {
		return false
	}

	titleHistory := make([]byte, 0)
	descriptionHistory := make([]byte, 0)

	if iter.err = iter.rows.Scan(
		&discussion.ID,
		&discussion.CreatedAt,
		&discussion.UpdatedAt,
		&discussion.DeletedAt,
		&discussion.Title,
		&discussion.AnonymityType,
		&discussion.ModeratorID,
		&discussion.IconURL,
		&discussion.Description,
		&titleHistory,
		&descriptionHistory,
		&discussion.DiscussionJoinability,
		&discussion.LastPostID,
		&discussion.LastPostCreatedAt,
		&discussion.ShuffleCount,
		&discussion.LockStatus,
	); iter.err != nil {
		logrus.WithError(iter.err).Error("iterator failed to scan row")
		return false
	}

	discussion.TitleHistory.RawMessage = titleHistory
	discussion.DescriptionHistory.RawMessage = descriptionHistory

	return true
}

func (iter *discussionIter) Close() error {
	if err := iter.err; err != nil {
		logrus.WithError(err).Error("iter error on close")
		return err
	}
	if err := iter.rows.Close(); err != nil {
		logrus.WithError(err).Error("iter rows close on close")
		return err
	}

	return nil
}
