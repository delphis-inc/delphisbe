package datastore

import sql2 "database/sql"

// Prepared Statements
type dbPrepStmts struct {
	// Post
	getPostByIDStmt                           *sql2.Stmt
	getPostsByDiscussionIDStmt                *sql2.Stmt
	getLastPostByDiscussionIDStmt             *sql2.Stmt
	getPostsByDiscussionIDFromCursorStmt      *sql2.Stmt
	putPostStmt                               *sql2.Stmt
	deletePostByIDStmt                        *sql2.Stmt
	deletePostByParticipantIDDiscussionIDStmt *sql2.Stmt

	// PostContents
	putPostContentsStmt *sql2.Stmt

	// Activity
	putActivityStmt *sql2.Stmt

	// Media
	putMediaRecordStmt *sql2.Stmt
	getMediaRecordStmt *sql2.Stmt

	// Discussion
	getDiscussionByLinkSlugStmt *sql2.Stmt

	// Discussion Archives
	getDiscussionArchiveByDiscussionIDStmt *sql2.Stmt
	upsertDiscussionArchiveStmt            *sql2.Stmt

	// Moderator
	getModeratorByUserIDStmt                *sql2.Stmt
	getModeratorByDiscussionIDStmt          *sql2.Stmt
	getModeratorByUserIDAndDiscussionIDStmt *sql2.Stmt
	getModeratedDiscussionsByUserIDStmt     *sql2.Stmt

	// Discussion Access
	getDiscussionsByUserAccessStmt       *sql2.Stmt
	getDiscussionUserAccessStmt          *sql2.Stmt
	getDUAForEverythingNotificationsStmt *sql2.Stmt
	getDUAForMentionNotificationsStmt    *sql2.Stmt
	upsertDiscussionUserAccessStmt       *sql2.Stmt
	deleteDiscussionUserAccessStmt       *sql2.Stmt

	// Requests
	getDiscussionRequestAccessByIDStmt         *sql2.Stmt
	getDiscussionAccessRequestsStmt            *sql2.Stmt
	getDiscussionAccessRequestByUserIDStmt     *sql2.Stmt
	getSentDiscussionAccessRequestsForUserStmt *sql2.Stmt
	putDiscussionAccessRequestStmt             *sql2.Stmt
	updateDiscussionAccessRequestStmt          *sql2.Stmt

	// AccessLinks
	getAccessLinkBySlugStmt           *sql2.Stmt
	getAccessLinkByDiscussionIDString *sql2.Stmt
	putAccessLinkForDiscussionString  *sql2.Stmt

	// DiscussionShuffleTimes
	getNextShuffleTimeForDiscussionIDString *sql2.Stmt
	putNextShuffleTimeForDiscussionIDString *sql2.Stmt
	getDiscussionsToShuffle                 *sql2.Stmt
	incrDiscussionShuffleCount              *sql2.Stmt

	// Viewers
	getViewerForDiscussionIDUserID *sql2.Stmt
	updateViewerLastViewed         *sql2.Stmt
}

const getPostByIDString = `
		SELECT p.id,
			p.created_at,
			p.updated_at,
			p.deleted_at,
			p.deleted_reason_code,
			p.discussion_id,
			p.participant_id,
			p.quoted_post_id,
			p.media_id,
			p.post_type,
			pc.id,
			pc.content,
			pc.mentioned_entities
		FROM posts p
		INNER JOIN post_contents pc
		ON p.post_content_id = pc.id
		WHERE p.id = $1;`

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
			p.post_type,
			pc.id,
			pc.content,
			pc.mentioned_entities
		FROM posts p
		INNER JOIN post_contents pc
		ON p.post_content_id = pc.id
		WHERE p.discussion_id = $1;`

const getPostsByDiscussionIDFromCursorString = `
		SELECT p.id,
			p.created_at,
			p.updated_at,
			p.deleted_at,
			p.deleted_reason_code,
			p.discussion_id,
			p.participant_id,
			p.quoted_post_id,
			p.media_id,
			p.post_type,
			pc.id,
			pc.content,
			pc.mentioned_entities
		FROM posts p
		INNER JOIN post_contents pc
		ON p.post_content_id = pc.id
		WHERE p.discussion_id = $1
		AND p.created_at < $2
		ORDER BY p.created_at desc
		LIMIT $3;`

const deletePostByIDString = `
		UPDATE posts
		SET deleted_at = now(),
			deleted_reason_code = $2,
			quoted_post_id = null,
			media_id = null
		WHERE id = $1
		RETURNING 
			id,
			created_at,
			updated_at,
			deleted_at,
			deleted_reason_code,
			discussion_id,
			participant_id,
			post_type;
`

const deletePostByParticipantIDDiscussionIDString = `
		UPDATE posts
		SET deleted_at = now(),
			deleted_reason_code = $3,
			quoted_post_id = null,
			media_id = null
		WHERE discussion_id = $1 AND
			participant_id = $2
		RETURNING id;
`

const getLastPostByDiscussionIDStmt = `
		SELECT p.id,
			p.created_at,
			p.updated_at,
			p.deleted_at,
			p.deleted_reason_code,
			p.discussion_id,
			p.participant_id,
			p.quoted_post_id,
			p.media_id,
			p.post_type,
			pc.id,
			pc.content,
			pc.mentioned_entities
		FROM posts p
		INNER JOIN post_contents pc
		ON p.post_content_id = pc.id
		WHERE p.discussion_id = $1
		ORDER BY p.created_at desc
		LIMIT 1;`

const putPostString = `
		INSERT INTO posts (
			id,
			discussion_id,
			participant_id,
			post_content_id,
			quoted_post_id,
			media_id,
			post_type
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING
			id,
			created_at,
			updated_at,
			discussion_id,
			participant_id,
			post_content_id,
			quoted_post_id,
			media_id,
			post_type;`

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

const getDiscussionByLinkSlugString = `
		SELECT d.id,
			d.created_at,
			d.updated_at,
			d.deleted_at,
			d.title,
			d.anonymity_type,
			d.moderator_id,
			d.icon_url,
			d.description,
			d.title_history,
			d.description_history,
			d.discussion_joinability,
			d.last_post_id,
			d.last_post_created_at,
			d.shuffle_count,
			d.lock_status
		FROM discussion_access_link dal
		INNER JOIN discussions d
		ON dal.discussion_id = d.id
		WHERE dal.link_slug = $1
			AND d.lock_status = false;`

// Discussion Archive
const getDiscussionArchiveByDiscussionIDString = `
		SELECT discussion_id,
			archived,
			created_at
		FROM discussion_archives
		WHERE discussion_id = $1;
`

const upsertDiscussionArchiveString = `
		INSERT INTO discussion_archives (
			discussion_id,
			archived
		) VALUES ($1, $2)
		ON CONFLICT (discussion_id)
		DO UPDATE SET created_at = now(),
			archived =$2
		RETURNING
			discussion_id,
			archived,
			created_at;`

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
		INNER JOIN discussions d
		ON m.id = d.moderator_id
		WHERE u.user_id = $1 LIMIT 1;`

const getModeratorByDiscussionIDString = `
		SELECT m.id,
			m.created_at,
			m.updated_at,
			m.deleted_at,
			m.user_profile_id
		FROM moderators m
		INNER JOIN discussions d
		ON m.id = d.moderator_id
		WHERE d.id = $1;
`

const getModeratorByUserIDAndDiscussionIDString = `
		SELECT m.id,
			m.created_at,
			m.updated_at,
			m.deleted_at,
			m.user_profile_id
		FROM moderators m
		INNER JOIN user_profiles u
		ON m.user_profile_id = u.id
		INNER JOIN discussions d
		ON m.id = d.moderator_id
		WHERE u.user_id = $1 and d.id = $2;`

const getModeratedDiscussionsByUserIDString = `
		SELECT d.id,
			d.created_at,
			d.updated_at,
			d.deleted_at,
			d.title,
			d.anonymity_type,
			d.moderator_id,
			d.icon_url,
			d.description,
			d.title_history,
			d.description_history,
			d.discussion_joinability,
			d.last_post_id,
			d.last_post_created_at,
			d.shuffle_count,
			d.lock_status
		FROM moderators m
		INNER JOIN user_profiles u
		ON m.user_profile_id = u.id
		INNER JOIN discussions d
		ON m.id = d.moderator_id
		WHERE u.user_id = $1;
`

// Discussion Access
const getDiscussionsByUserAccessString = `
		SELECT d.id,
			d.created_at,
			d.updated_at,
			d.deleted_at,
			d.title,
			d.anonymity_type,
			d.moderator_id,
			d.icon_url,
			d.description,
			d.title_history,
			d.description_history,
			d.discussion_joinability,
			d.last_post_id,
			d.last_post_created_at,
			d.shuffle_count,
			d.lock_status
		FROM discussion_user_access dua
		INNER JOIN discussions d
			ON dua.discussion_id = d.id
		WHERE dua.user_id = $1
			AND dua.state = $2
			AND d.deleted_at is null
		ORDER BY d.last_post_created_at desc;`

const getDiscussionUserAccessString = `
		SELECT 	discussion_id,
			user_id,
			state,
			request_id,
			notif_setting,
			created_at,
			updated_at,
			deleted_at
		FROM discussion_user_access
		WHERE discussion_id = $1
			AND user_id = $2;`

const getDUAForEverythingNotificationsString = `
		SELECT 	discussion_id,
			user_id,
			state,
			request_id,
			notif_setting,
			created_at,
			updated_at,
			deleted_at
		FROM discussion_user_access
		WHERE discussion_id = $1
			AND user_id != $2
			AND state = 'ACTIVE'
			AND notif_setting = 'EVERYTHING';`

const getDUAForMentionNotificationsString = `
		SELECT 	discussion_id,
			user_id,
			state,
			request_id,
			notif_setting,
			created_at,
			updated_at,
			deleted_at
		FROM discussion_user_access
		WHERE discussion_id = $1
			AND user_id != $2
			AND user_id = ANY($3)
			AND state = 'ACTIVE'
			AND notif_setting = 'MENTIONS';` // We could also check if notif_setting != NONE if we wanted to treat these notifs differently

const upsertDiscussionUserAccessString = `
		INSERT INTO discussion_user_access (
			discussion_id,
			user_id,
			state,
			request_id,
			notif_setting
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (discussion_id, user_id)
		DO UPDATE SET state = $3,
			request_id = $4,
			notif_setting = $5
		RETURNING
			discussion_id,
			user_id,
			state,
			request_id,
			notif_setting,
			created_at,
			updated_at,
			deleted_at;`

const deleteDiscussionUserAccessString = `
		UPDATE discussion_user_access
		SET deleted_at = now()
		WHERE discussion_id = $1
			AND user_id = $2
		RETURNING
			discussion_id,
			user_id,
			created_at,
			updated_at,
			deleted_at;`

const getDiscussionRequestAccessByIDString = `
		SELECT id,
			user_id,
			discussion_id,
			created_at,
			updated_at,
			status
		FROM discussion_user_requests
		WHERE id = $1;`

const getDiscussionAccessRequestsString = `
		SELECT id,
			user_id,
			discussion_id,
			created_at,
			updated_at,
			status
		FROM discussion_user_requests
		WHERE discussion_id = $1
			AND status = 'PENDING'
			AND deleted_at is null;`

const getDiscussionAccessRequestByUserIDString = `
		SELECT id,
			user_id,
			discussion_id,
			created_at,
			updated_at,
			status
		FROM discussion_user_requests
		WHERE discussion_id = $1
		    AND user_id = $2
			AND deleted_at is null;`

const getSentDiscussionAccessRequestsForUserString = `
		SELECT id,
			user_id,
			discussion_id,
			created_at,
			updated_at,
			status
		FROM discussion_user_requests
		WHERE user_id = $1
			AND deleted_at is null;`

const putDiscussionAccessRequestString = `
		INSERT INTO discussion_user_requests (
			id,
			user_id,
			discussion_id,
			status
		) VALUES ($1, $2, $3, $4)
		RETURNING id,
			user_id,
			discussion_id,
			created_at,
			updated_at,
			status;`

const updateDiscussionAccessRequestString = `
		UPDATE discussion_user_requests
		SET status = $2
		WHERE id = $1
		RETURNING id,
			user_id,
			discussion_id,
			created_at,
			updated_at,
			status;`

// AccessLinks
const getAccessLinkBySlugString = `
		SELECT discussion_id,
			link_slug,
			created_at,
			updated_at,
			deleted_at
		FROM discussion_access_link
		WHERE link_slug = $1
			AND deleted_at is null;`

const getAccessLinkByDiscussionIDString = `
		SELECT discussion_id,
			link_slug,
			created_at,
			updated_at,
			deleted_at
		FROM discussion_access_link
		WHERE discussion_id = $1
			AND deleted_at is null LIMIT 1;`

const putAccessLinkForDiscussionString = `
		INSERT into discussion_access_link (
			discussion_id,
			link_slug
		) VALUES ($1, $2)
		RETURNING discussion_id,
			link_slug,
			created_at,
			updated_at,
			deleted_at;`

const getNextShuffleTimeForDiscussionIDString = `
		SELECT discussion_id,
			shuffle_time
		from discussion_shuffle_time
		WHERE discussion_id = $1;`

const putNextShuffleTimeForDiscussionIDString = `
		INSERT INTO discussion_shuffle_time (
			discussion_id,
			shuffle_time
		) VALUES ($1, $2)
		ON CONFLICT (discussion_id)
		DO UPDATE SET shuffle_time = $2
		RETURNING
			discussion_id,
			shuffle_time;`

const getDiscussionsToShuffle = `
		SELECT d.id,
			d.shuffle_count
		FROM discussion_shuffle_time s
		JOIN discussions d ON d.id = s.discussion_id
		WHERE shuffle_time is not NULL 
		AND shuffle_time <= $1
		AND d.lock_status = false;`

// This may cause multiple updates to happen to the same row but since
// shuffling is sort of idempotent (no expected outcome) it's a good
// non-locking approach for now!
const incrDiscussionShuffleCount = `
		UPDATE discussions
		SET shuffle_count = shuffle_count + 1
		WHERE id = $1
		RETURNING
			shuffle_count;`

const getViewerForDiscussionIDUserID = `
		SELECT
			id,
			created_at,
			updated_at,
			last_viewed,
			last_viewed_post_id,
			discussion_id,
			user_id
		FROM viewers
		WHERE 
			discussion_id = $1 
			AND user_id = $2
			AND deleted_at is NULL;`

const updateViewerLastViewed = `
		UPDATE viewers
		SET last_viewed = $2, last_viewed_post_id = $3
		WHERE id = $1
		RETURNING 
			id,
			created_at,
			updated_at,
			last_viewed,
			last_viewed_post_id,
			discussion_id,
			user_id;`
