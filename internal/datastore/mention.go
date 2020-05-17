package datastore

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pkg/errors"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/sirupsen/logrus"
)

func (d *delphisDB) PutActivity(ctx context.Context, tx *sql.Tx, post *model.Post) error {
	logrus.Infof("PutActivity::SQL Create")
	if err := d.initializeStatements(ctx); err != nil {
		logrus.WithError(err).Error("PutActivity::failed to initialize statements")
		return err
	}

	for _, entityID := range post.PostContent.MentionedEntities {
		s := strings.Split(entityID, ":")

		// Don't record mentions where a user tags themselves. This can also be handled on the frontend
		if s[1] != *post.ParticipantID {
			_, err := tx.StmtContext(ctx, d.prepStmts.putActivityStmt).ExecContext(
				ctx,
				post.ParticipantID,
				post.PostContent.ID,
				s[0],
				s[1],
			)
			if err != nil {
				logrus.WithError(err).Error("failed to execute putActivityStmt")
				return errors.Wrap(err, "failed to putMention")
			}
		}
	}
	return nil
}
