package datastore

import (
	"context"
	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/nedrocks/delphisbe/graph/model"
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
	moderator := model.Moderator{}
	if err := d.sql.Preload("Moderator").First(&moderator, model.Moderator{ID: moderatorID}).Related(&discussion).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		logrus.WithError(err).Errorf("GetDiscussionByModeratorID::Failed getting discussion by moderator ID")
		return nil, err
	}

	return &discussion, nil
}

func (d *delphisDB) GetDiscussionsAutoPost(ctx context.Context) AutoPostDiscussionIter {
	logrus.Debug("GetDiscussionsAutoPost::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetDiscussionsAutoPost::failed to initialize statements")
		return &autoPostDiscussionIter{err: err}
	}

	rows, err := d.prepStmts.getDiscussionsForAutoPostStmt.QueryContext(
		ctx,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to query GetDiscussionsAutoPost")
		return &autoPostDiscussionIter{err: err}
	}

	return &autoPostDiscussionIter{
		ctx:  ctx,
		rows: rows,
	}
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
			"Title":         discussion.Title,
			"AnonymityType": discussion.AnonymityType,
			"AutoPost":      discussion.AutoPost,
			"IdleMinutes":   discussion.IdleMinutes,
			"PublicAccess":  discussion.PublicAccess,
		}).First(&found).Error; err != nil {
			logrus.WithError(err).Errorf("UpsertDiscussion::Failed updating disucssion object")
			return nil, err
		}
	}
	return &found, nil
}

func (d *delphisDB) GetDiscussionTags(ctx context.Context, id string) TagIter {
	logrus.Debug("GetDiscussionTags::SQL Query")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("GetDiscussionTags::failed to initialize statements")
		return &tagIter{err: err}
	}

	rows, err := d.prepStmts.getDiscussionTagsStmt.QueryContext(
		ctx,
		id,
	)
	if err != nil {
		logrus.WithError(err).Error("failed to query GetDiscussionTagsStmt")
		return &tagIter{err: err}
	}

	return &tagIter{
		ctx:  ctx,
		rows: rows,
	}
}

func (d *delphisDB) PutDiscussionTags(ctx context.Context, tx *sql.Tx, tag model.Tag) (*model.Tag, error) {
	logrus.Debug("PutDiscussionTags::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("PutDiscussionTags::failed to initialize statements")
		return nil, err
	}

	if err := tx.StmtContext(ctx, d.prepStmts.putDiscussionTagsStmt).QueryRowContext(
		ctx,
		tag.ID,
		tag.Tag,
	).Scan(
		&tag.ID,
		&tag.Tag,
		&tag.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return &model.Tag{}, nil
		}
		logrus.WithError(err).Error("failed to execute putDiscussionTagsStmt")
		return nil, err
	}

	return &tag, nil
}

func (d *delphisDB) DeleteDiscussionTags(ctx context.Context, tx *sql.Tx, tag model.Tag) (*model.Tag, error) {
	logrus.Debug("DeleteDiscussionTags::SQL Delete")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("DeleteDiscussionTags::failed to initialize statements")
		return nil, err
	}

	if err := tx.StmtContext(ctx, d.prepStmts.deleteDiscussionTagsStmt).QueryRowContext(
		ctx,
		tag.ID,
		tag.Tag,
	).Scan(
		&tag.ID,
		&tag.Tag,
		&tag.CreatedAt,
		&tag.DeletedAt,
	); err != nil {
		logrus.WithError(err).Error("failed to execute deleteDiscussionTags")
		return nil, err
	}

	return &tag, nil
}

type autoPostDiscussionIter struct {
	err  error
	ctx  context.Context
	rows *sql.Rows
}

func (iter *autoPostDiscussionIter) Next(discussion *model.DiscussionAutoPost) bool {
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

	if iter.err = iter.rows.Scan(
		&discussion.ID,
		&discussion.IdleMinutes,
	); iter.err != nil {
		logrus.WithError(iter.err).Error("iterator failed to scan row")
		return false
	}

	return true
}

func (iter *autoPostDiscussionIter) Close() error {
	if err := iter.err; err != nil {
		logrus.WithError(err).Error("iter error on close")
		return err
	}
	if err := iter.rows.Close(); err != nil {
		logrus.WithError(err).Error("iter rows close on close")
	}

	return nil
}
