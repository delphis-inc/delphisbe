package datastore

import (
	"context"

	"go.uber.org/multierr"

	"github.com/nedrocks/delphisbe/graph/model"

	"github.com/sirupsen/logrus"
)

type TestData struct {
	Discussions  []model.Discussion
	Participants []model.Participant
	Posts        []model.Post
}

func (d *delphisDB) CreateTestTables(ctx context.Context, data TestData) (func() error, error) {
	if err := d.createTestTables(ctx); err != nil {
		logrus.WithError(err).Error("failed to create test tables")
		return nil, nil
	}

	if err := d.populateTestTables(ctx, data); err != nil {
		logrus.WithError(err).Error("failed to populate test tables")
		return nil, err
	}

	return d.close, nil
}

func (d *delphisDB) createTestTables(ctx context.Context) error {
	d.readyMu.Lock()
	defer d.readyMu.Unlock()

	createTableQueries := []string{
		`CREATE TABLE IF NOT EXISTS posts (
			id varchar(36) PRIMARY KEY,
			created_at timestamp with time zone default current_timestamp not null,
			updated_at timestamp with time zone default current_timestamp not null,
			deleted_at timestamp with time zone,
			deleted_reason_code varchar(16),
			discussion_id varchar(36),
			participant_id varchar(36),
			post_content_id varchar(36) not null,
			quoted_post_id varchar(36),
			media_id varchar(36),
			imported_content_id varchar(36)
		);`,

		`CREATE TABLE IF NOT EXISTS post_contents (
			id varchar(36) PRIMARY KEY,
			content text,
			created_at timestamp with time zone default current_timestamp not null,
			updated_at timestamp with time zone default current_timestamp not null,
			mentioned_entities varchar(50)[]
		);`,

		`CREATE TABLE IF NOT EXISTS discussions (
			id varchar(36) PRIMARY KEY,
			created_at timestamp with time zone default current_timestamp not null,
			updated_at timestamp with time zone default current_timestamp not null,
			deleted_at timestamp with time zone,
			title varchar(256) not null,
			anonymity_type varchar(36) not null,
			moderator_id varchar(36),
			auto_post boolean default true not null,
			idle_minutes int default 300 not null
		);`,

		`CREATE TABLE IF NOT EXISTS participants (
			id varchar(36) PRIMARY KEY,
			participant_id int not null,
			created_at timestamp with time zone default current_timestamp not null,
			updated_at timestamp with time zone default current_timestamp not null,
			deleted_at timestamp with time zone,
			discussion_id varchar(36) not null,
			viewer_id varchar(36) not null,
			user_id varchar(36),
			flair_id varchar(36),
		    is_anonymous boolean NOT NULL DEFAULT True,
			gradient_color varchar(36),
			has_joined boolean NOT NULL DEFAULT FALSE
		);`,

		`CREATE TABLE IF NOT EXISTS media (
			id varchar(36) not null PRIMARY KEY,
			created_at timestamp with time zone default current_timestamp not null,
			deleted_at timestamp with time zone,
			deleted_reason_code varchar(16),
			media_type varchar(16),
			media_size json not null
		);`,

		`CREATE TABLE IF NOT EXISTS activity (
			participant_id varchar(36) not null,
			post_content_id varchar(36) not null,
			entity_id varchar(36) not null,
			entity_type varchar(36) not null,
			created_at timestamp with time zone default current_timestamp not null,
			PRIMARY KEY(participant_id, entity_id, created_at)
		);`,

		`CREATE TABLE IF NOT EXISTS moderators (
			id varchar(36) PRIMARY KEY,
			created_at timestamp with time zone not null,
			updated_at timestamp with time zone not null,
			deleted_at timestamp with time zone,
			user_profile_id varchar(36) not null
		);`,

		`CREATE TABLE IF NOT EXISTS user_profiles (
			id varchar(36) PRIMARY KEY,
			created_at timestamp with time zone not null,
			updated_at timestamp with time zone not null,
			deleted_at timestamp with time zone,
			display_name varchar(128) not null,
			user_id varchar(36),
			twitter_handle varchar(36)
		);`,

		`CREATE TABLE IF NOT EXISTS imported_contents (
			id varchar(36) PRIMARY KEY,
			created_at timestamp with time zone default current_timestamp not null,
			content_name text not null,
			content_type text not null,
			link text not null,
			overview text not null,
			source text not null
		);`,

		`CREATE TABLE IF NOT EXISTS discussion_tags (
			discussion_id varchar(36) not null,
			deleted_at timestamp with time zone,
			tag varchar(40) not null,
			created_at timestamp with time zone default current_timestamp not null,
			PRIMARY KEY(discussion_id, tag)
		);`,

		`CREATE TABLE IF NOT EXISTS imported_content_tags (
			imported_content_id varchar(36) not null,
			deleted_at timestamp with time zone,
			tag varchar(40) not null,
			created_at timestamp with time zone default current_timestamp not null,
			PRIMARY KEY(imported_content_id, tag)
		);`,

		`CREATE TABLE IF NOT EXISTS discussion_ic_queue (
			discussion_id varchar(36) not null,
			imported_content_id varchar(36) not null,
			created_at timestamp with time zone default current_timestamp not null,
			updated_at timestamp with time zone default current_timestamp not null,
			deleted_at timestamp with time zone,
			posted_at timestamp with time zone,
			matching_tags varchar(40)[]
		);`,

		// Add foreign keys
		`ALTER TABLE posts
			ADD CONSTRAINT posts_discussions_fk_c34cae6d6fc5 FOREIGN KEY (discussion_id) REFERENCES discussions (id) MATCH FULL,
			ADD CONSTRAINT posts_participants_fk_c94a4fb2438b FOREIGN KEY (participant_id) REFERENCES participants (id) MATCH FULL,
			ADD CONSTRAINT posts_post_contents_fk_777ecc8c7969 FOREIGN KEY (post_content_id) REFERENCES post_contents (id) MATCH FULL,
			ADD CONSTRAINT posts_quoted_post_id_eaa15cd7531b FOREIGN KEY (quoted_post_id) REFERENCES posts (id) MATCH FULL,
			ADD CONSTRAINT posts_media_fk_b783eacfac89 FOREIGN KEY (media_id) REFERENCES media (id) MATCH FULL;`,
	}

	tx, err := d.pg.BeginTx(ctx, nil)
	if err != nil {
		logrus.WithError(err).Error("failed to start CreateTestTables transaction")
		return err
	}

	for _, query := range createTableQueries {
		if _, err = tx.ExecContext(ctx, query); err != nil {
			logrus.WithError(err).Errorf("failed table creation for %v\n", query)
			return tx.Rollback()
		}
	}

	return tx.Commit()
}

func (d *delphisDB) populateTestTables(ctx context.Context, data TestData) error {
	if err := d.writeDiscussions(ctx, data.Discussions); err != nil {
		logrus.WithError(err).Error("failed to write discussions to test table")
		return err
	}

	if err := d.writeParticipants(ctx, data.Participants); err != nil {
		logrus.WithError(err).Error("failed to write participants to test table")
		return err
	}

	if err := d.writePostsAndContents(ctx, data.Posts); err != nil {
		logrus.WithError(err).Error("failed to write post and contents to test table")
		return err
	}

	return nil
}

func (d *delphisDB) writeDiscussions(ctx context.Context, testDiscussions []model.Discussion) error {
	// Iterate over test data to create test records
	for _, discussion := range testDiscussions {
		logrus.Debugf("In here for Disc: %+v\n", discussion)
		if _, err := d.UpsertDiscussion(ctx, discussion); err != nil {
			logrus.WithError(err).Error("failed to upsert test discussion")
			return err
		}
	}

	return nil
}

func (d *delphisDB) writeParticipants(ctx context.Context, testParticipants []model.Participant) error {
	for _, participant := range testParticipants {
		logrus.Debugf("In here for participant: %+v\n", participant)
		if _, err := d.UpsertParticipant(ctx, participant); err != nil {
			logrus.WithError(err).Error("failed to upsert test participant")
			return err
		}
	}
	return nil
}

func (d *delphisDB) writePostsAndContents(ctx context.Context, testPosts []model.Post) error {
	tx, err := d.BeginTx(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to create tx")
		return err
	}
	for _, post := range testPosts {
		logrus.Debugf("In here for posts: %+v\n", post)
		if err := d.PutPostContent(ctx, tx, *post.PostContent); err != nil {
			logrus.WithError(err).Error("failed to put test post contents")

			// Rollback on errors
			if txErr := d.RollbackTx(ctx, tx); txErr != nil {
				logrus.WithError(txErr).Error("failed to rollback tx")
				return multierr.Append(err, txErr)
			}
			return err
		}

		if _, err := d.PutPost(ctx, tx, post); err != nil {
			logrus.WithError(err).Error("failed to put test post")

			// Rollback on errors
			if txErr := d.RollbackTx(ctx, tx); txErr != nil {
				logrus.WithError(txErr).Error("failed to rollback tx")
				return multierr.Append(err, txErr)
			}
			return err
		}
	}

	return d.CommitTx(ctx, tx)
}

func (d *delphisDB) close() error {
	if err := d.sql.Close(); err != nil {
		return err
	}
	if err := d.pg.Close(); err != nil {
		return err
	}
	return nil
}
