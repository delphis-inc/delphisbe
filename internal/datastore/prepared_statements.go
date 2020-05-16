package datastore

import sql2 "database/sql"

// Prepared Statements
type dbPrepStmts struct {
	// Post
	putPostStmt                *sql2.Stmt
	getPostsByDiscussionIDStmt *sql2.Stmt

	// PostContents
	putPostContentsStmt *sql2.Stmt
}

const putPostString = `
		INSERT INTO posts (
			id,
			discussion_id,
			participant_id,
			post_content_id,
			quoted_post_id
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			created_at,
			updated_at,
			discussion_id,
			participant_id,
			post_content_id,
			quoted_post_id;`
