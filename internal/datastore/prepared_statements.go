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
	getDiscussionsForAutoPostStmt *sql2.Stmt

	// Moderator
	getModeratorByUserIDStmt                *sql2.Stmt
	getModeratorByUserIDAndDiscussionIDStmt *sql2.Stmt

	// ImportedContent
	getImportedContentByIDStmt                    *sql2.Stmt
	getImportedContentForDiscussionStmt           *sql2.Stmt
	getScheduledImportedContentByDiscussionIDStmt *sql2.Stmt
	putImportedContentStmt                        *sql2.Stmt
	putImportedContentDiscussionQueueStmt         *sql2.Stmt
	updateImportedContentDiscussionQueueStmt      *sql2.Stmt

	// Tags
	getImportedContentTagsStmt *sql2.Stmt
	getDiscussionTagsStmt      *sql2.Stmt
	getMatchingTagsStmt        *sql2.Stmt
	putImportedContentTagsStmt *sql2.Stmt
	putDiscussionTagsStmt      *sql2.Stmt
	deleteDiscussionTagsStmt   *sql2.Stmt

	// Discussion Access
	getDiscussionsByUserAccessStmt *sql2.Stmt
	getDiscussionUserAccessStmt    *sql2.Stmt
	upsertDiscussionUserAccessStmt *sql2.Stmt
	deleteDiscussionUserAccessStmt *sql2.Stmt

	// InvitesRequests
	getDiscussionInviteByIDStmt                            *sql2.Stmt
	getDiscussionRequestAccessByIDStmt                     *sql2.Stmt
	getDiscussionInvitesForUserStmt                        *sql2.Stmt
	getSentDiscussionInvitesForUserStmt                    *sql2.Stmt
	getDiscussionAccessRequestsStmt                        *sql2.Stmt
	getDiscussionAccessRequestByUserIDStmt                 *sql2.Stmt
	getSentDiscussionAccessRequestsForUserStmt             *sql2.Stmt
	putDiscussionInviteRecordStmt                          *sql2.Stmt
	putDiscussionAccessRequestStmt                         *sql2.Stmt
	updateDiscussionInviteRecordStmt                       *sql2.Stmt
	updateDiscussionAccessRequestStmt                      *sql2.Stmt
	getInvitedTwitterHandlesByDiscussionIDAndInviterIDStmt *sql2.Stmt

	// AccessLinks
	getAccessLinkBySlugStmt           *sql2.Stmt
	getAccessLinkByDiscussionIDString *sql2.Stmt
	putAccessLinkForDiscussionString  *sql2.Stmt

	// DiscussionShuffleTimes
	getNextShuffleTimeForDiscussionIDString *sql2.Stmt
	putNextShuffleTimeForDiscussionIDString *sql2.Stmt
	getDiscussionsToShuffle                 *sql2.Stmt
	incrDiscussionShuffleCount              *sql2.Stmt
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
			p.imported_content_id,
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
			p.imported_content_id,
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
			p.imported_content_id,
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
			media_id = null,
			imported_content_id = null
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
			media_id = null,
			imported_content_id = null
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
			p.imported_content_id,
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
			imported_content_id,
			post_type
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			id,
			created_at,
			updated_at,
			discussion_id,
			participant_id,
			post_content_id,
			quoted_post_id,
			media_id,
			imported_content_id,
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

const getDiscussionsForAutoPostString = `
		SELECT id,
			idle_minutes
		FROM discussions
		WHERE auto_post = true`

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

const getImportedContentByIDString = `
		SELECT id,
			created_at,
			content_name,
			content_type,
			link,
			overview,
			source
		FROM imported_contents
		WHERE id = $1;`

// TODO: We may just want to write these records into the DB or cache with a TTL.
// Subquery gets the imported_contents that matches up with a discussions tag.
// It then checks the results against the imported_ic_queue table to see what has
// been posted or is scheduled.
// Finally, we join with imported_contents to retrieve the data.
const getImportedContentForDiscussionString = `
		SELECT i.id,
			i.created_at,
			i.content_name,
			i.content_type,
			i.link,
			i.overview,
			i.source,
			d.matching_tags
		FROM (
			SELECT d.discussion_id,
				i.imported_content_id,
				array_agg(i.tag) matching_tags
			FROM discussion_tags d
			INNER JOIN imported_content_tags i
			ON d.tag = i.tag
			WHERE d.discussion_id = $1
				AND NOT EXISTS (
					SELECT
					FROM discussion_ic_queue q
					WHERE d.discussion_id = q.discussion_id
						AND i.imported_content_id = q.imported_content_id
			)
			GROUP BY d.discussion_id, i.imported_content_id
		) d
		INNER JOIN imported_contents i
		ON i.id = d.imported_content_id
		ORDER BY i.created_at desc
		LIMIT $2;`

// TODO: Do we want to limit this to schedule articles in the past 24 or 48 hours so we don't post old stories?
// TODO: What do we want to order by? When the article was posted? When the article was added to the queue?
// Subquery gets the unposted imported contents for a discussion as these are scheduled
// Then, we get the imported contents data from the table
const getScheduledImportedContentByDiscussionIDString = `
		SELECT i.id,
			i.created_at,
			i.content_name,
			i.content_type,
			i.link,
			i.overview,
			i.source,
			q.matching_tags
		FROM (
			SELECT imported_content_id,
			matching_tags,
			updated_at
			FROM discussion_ic_queue
			WHERE discussion_id = $1
				AND posted_at is null
				AND updated_at > now() - interval '48 hours'
		) q
		INNER JOIN imported_contents i
		ON i.id = q.imported_content_id
		ORDER BY q.updated_at desc;`

const putImportedContentString = `
		INSERT INTO imported_contents (
			id,
			content_name,
			content_type,
			link,
			overview,
			source
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			created_at,
			content_name,
			content_type,
			link,
			overview,
			source;`

const putImportedContentDiscussionQueueString = `
		INSERT INTO discussion_ic_queue (
			discussion_id,
			imported_content_id,
			posted_at,
			matching_tags
		) VALUES ($1, $2, $3, $4)
		RETURNING
			discussion_id,
			imported_content_id,
			created_at,
			updated_at,
			deleted_at,
			posted_at,
			matching_tags;`

const updateImportedContentDiscussionQueueString = `
		UPDATE discussion_ic_queue
		SET posted_at = $3
		WHERE discussion_id = $1
			AND imported_content_id = $2
			AND posted_at is null
		RETURNING
			discussion_id,
			imported_content_id,
			created_at,
			updated_at,
			deleted_at,
			posted_at,
			matching_tags;`

const getImportedContentTagsString = `
		SELECT imported_content_id
			tag,
			created_at
		FROM imported_content_tags
		WHERE imported_content_id = $1
			AND deleted_at is null
		ORDER BY created_at desc;`

const getDiscussionTagsString = `
		SELECT discussion_id,
			tag,
			created_at
		FROM discussion_tags
		WHERE discussion_id = $1
			AND deleted_at is null
		ORDER BY created_at desc;`

const getMatchingTagsString = `
		SELECT array_agg(i.tag) matching_tags
		FROM discussion_tags d
		INNER JOIN imported_content_tags i
			ON i.tag = d.tag
		WHERE discussion_id = $1
			AND imported_content_id = $2
		GROUP BY discussion_id, imported_content_id;`

const putImportedContentTagsString = `
		INSERT INTO imported_content_tags (
			imported_content_id,
			tag
		) VALUES ($1, $2)
		RETURNING
			imported_content_id,
			tag,
			created_at;`

const putDiscussionTagsString = `
		INSERT INTO discussion_tags (
			discussion_id,
			tag
		) VALUES ($1, $2)
		ON CONFLICT (discussion_id, tag)
		DO UPDATE SET deleted_at = null
		RETURNING
			discussion_id,
			tag,
			created_at;`

const deleteDiscussionTagsString = `
		UPDATE discussion_tags
		SET deleted_at = now()
		WHERE discussion_id = $1
			AND tag = $2
		RETURNING
			discussion_id,
			tag,
			created_at,
			deleted_at;`

// Discussion Access
const getDiscussionsByUserAccessString = `
		SELECT d.id,
			d.created_at,
			d.updated_at,
			d.deleted_at,
			d.title,
			d.anonymity_type,
			d.moderator_id,
			d.auto_post,
			d.icon_url,
			d.idle_minutes,
			d.description,
			d.title_history,
			d.description_history,
			d.discussion_joinability,
			d.last_post_id,
			d.last_post_created_at,
			d.shuffle_count
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
			created_at,
			updated_at,
			deleted_at
		FROM discussion_user_access
		WHERE discussion_id = $1
			AND user_id = $2;`

const upsertDiscussionUserAccessString = `
		INSERT INTO discussion_user_access (
			discussion_id,
			user_id,
			state,
			request_id
		) VALUES ($1, $2, $3, $4)
		ON CONFLICT (discussion_id, user_id)
		DO UPDATE SET state = $3,
			request_id = $4
		RETURNING
			discussion_id,
			user_id,
			state,
			request_id,
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

// Invites and Requests
const getDiscussionInviteByIDString = `
		SELECT id,
			user_id,
			discussion_id,
			invite_from_participant_id,
			created_at,
			updated_at,
			status,
			invite_type
		FROM discussion_user_invitations
		WHERE id = $1;`

const getDiscussionRequestAccessByIDString = `
		SELECT id,
			user_id,
			discussion_id,
			created_at,
			updated_at,
			status
		FROM discussion_user_requests
		WHERE id = $1;`

const getDiscussionInvitesForUserString = `
		SELECT id,
			user_id,
			discussion_id,
			invite_from_participant_id,
			created_at,
			updated_at,
			status,
			invite_type
		FROM discussion_user_invitations
		WHERE user_id = $1
			AND deleted_at is null
			AND status = $2;`

const getSentDiscussionInvitesForUserString = `
		SELECT id,
			user_id,
			discussion_id,
			invite_from_participant_id,
			created_at,
			updated_at,
			status,
			invite_type
		FROM discussion_user_invitations
		WHERE invite_from_participant_id = $1
			AND deleted_at is null;`

const getDiscussionAccessRequestsString = `
		SELECT id,
			user_id,
			discussion_id,
			created_at,
			updated_at,
			status
		FROM discussion_user_requests
		WHERE discussion_id = $1
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

const putDiscussionInviteRecordString = `
		INSERT INTO discussion_user_invitations (
			id,
			user_id,
			discussion_id,
			invite_from_participant_id,
			status,
			invite_type
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id,
			user_id,
			discussion_id,
			invite_from_participant_id,
			created_at,
			updated_at,
			status,
			invite_type;`

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

const updateDiscussionInviteRecordString = `
		UPDATE discussion_user_invitations
		SET status = $2
		WHERE id = $1
		RETURNING id,
			user_id,
			discussion_id,
			invite_from_participant_id,
			created_at,
			updated_at,
			status,
			invite_type;`

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

const getInvitedTwitterHandlesByDiscussionIDAndInviterIDString = `
		SELECT up.twitter_handle
		FROM discussion_user_invitations dui
			LEFT JOIN users u ON u.id = dui.user_id
			LEFT JOIN user_profiles up ON up.user_id = u.id
		WHERE dui.discussion_id=$1 AND dui.invite_from_participant_id=$2;`

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
		AND shuffle_time <= $1;
`

// This may cause multiple updates to happen to the same row but since
// shuffling is sort of idempotent (no expected outcome) it's a good
// non-locking approach for now!
const incrDiscussionShuffleCount = `
		UPDATE discussions
		SET shuffle_count = shuffle_count + 1
		WHERE id = $1
		RETURNING
			shuffle_count;
`
