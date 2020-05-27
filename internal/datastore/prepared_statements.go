package datastore

import sql2 "database/sql"

// Prepared Statements
type dbPrepStmts struct {
	// Post
	putPostStmt                *sql2.Stmt
	getPostsByDiscussionIDStmt *sql2.Stmt

	// PostContents
	putPostContentsStmt *sql2.Stmt

	// Activity
	putActivityStmt *sql2.Stmt

	// Media
	putMediaRecordStmt *sql2.Stmt
	getMediaRecordStmt *sql2.Stmt

	// Moderator
	getModeratorByUserIDStmt *sql2.Stmt
}

const putPostString = `
		INSERT INTO posts (
			id,
			discussion_id,
			participant_id,
			post_content_id,
			quoted_post_id,
			media_id
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			created_at,
			updated_at,
			discussion_id,
			participant_id,
			post_content_id,
			quoted_post_id,
			media_id;`

const getPostsByDiscussionIDString = `
		SELECT p.id,
			p.created_at,
			p.updated_at,
			p.deleted_at,
			p.deleted_reason_code,
			p.discussion_id,
			p.participant_id,
			p.quoted_post_id,
			p.media_id,
			pc.id,
			pc.content
		FROM posts p
		INNER JOIN post_contents pc
		ON p.post_content_id = pc.id
		WHERE p.discussion_id = $1;`

const putPostContentsString = `
		INSERT INTO post_contents (
			id,
			content,
			mentioned_entities
		) VALUES ($1, $2, $3);`

const putActivityString = `
		INSERT INTO activity (
			participant_id,
			post_content_id,
			entity_id,
			entity_type
		) VALUES ($1, $2, $3, $4);`

const putMediaRecordString = `
		INSERT INTO media (
			id,
			media_type,
			media_size
		) VALUES ($1, $2, $3);`

const getMediaRecordString = `
		SELECT id,
			created_at,
			deleted_at,
			deleted_reason_code,
			media_type,
			media_size
		FROM media
		WHERE id = $1;`

// Currently only care if you are a mod, not checking on discussion mods
const getModeratorByUserIDString = `
		SELECT m.id,
			m.created_at,
			m.updated_at,
			m.deleted_at,
			m.user_profile_id
		FROM moderators m
		INNER JOIN user_profiles u
		ON m.user_profile_id = u.id
		WHERE u.user_id = $1 LIMIT 1;`
