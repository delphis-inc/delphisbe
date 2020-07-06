package datastore

import (
	"fmt"

	"github.com/DATA-DOG/go-sqlmock"
)

func mockPreparedStatements(mock sqlmock.Sqlmock) {
	// Prepare statements
	mock.ExpectPrepare(getPostByIDString)
	mock.ExpectPrepare(getPostsByDiscussionIDString)
	mock.ExpectPrepare(getLastPostByDiscussionIDStmt)
	mock.ExpectPrepare(getPostsByDiscussionIDFromCursorString)
	mock.ExpectPrepare(putPostString)
	mock.ExpectPrepare(deletePostByIDString)
	mock.ExpectPrepare(deletePostByParticipantIDDiscussionIDString)
	mock.ExpectPrepare(putPostContentsString)
	mock.ExpectPrepare(putActivityString)
	mock.ExpectPrepare(putMediaRecordString)
	mock.ExpectPrepare(getMediaRecordString)
	mock.ExpectPrepare(getDiscussionsForAutoPostString)
	mock.ExpectPrepare(getPublicDiscussionsString)
	mock.ExpectPrepare(getModeratorByUserIDString)
	mock.ExpectPrepare(getModeratorByUserIDAndDiscussionIDString)
	mock.ExpectPrepare(getImportedContentByIDString)
	mock.ExpectPrepare(getImportedContentForDiscussionString)
	mock.ExpectPrepare(getScheduledImportedContentByDiscussionIDString)
	mock.ExpectPrepare(putImportedContentString)
	mock.ExpectPrepare(putImportedContentDiscussionQueueString)
	mock.ExpectPrepare(updateImportedContentDiscussionQueueString)
	mock.ExpectPrepare(getImportedContentTagsString)
	mock.ExpectPrepare(getDiscussionTagsString)
	mock.ExpectPrepare(getMatchingTagsString)
	mock.ExpectPrepare(putImportedContentTagsString)
	mock.ExpectPrepare(putDiscussionTagsString)
	mock.ExpectPrepare(deleteDiscussionTagsString)
	mock.ExpectPrepare(getDiscussionsByFlairTemplateForUserString)
	mock.ExpectPrepare(getDiscussionsByUserAccessForUserString)
	mock.ExpectPrepare(getDiscussionFlairAccessString)
	mock.ExpectPrepare(upsertDiscussionFlairAccessString)
	mock.ExpectPrepare(upsertDiscussionUserAccessString)
	mock.ExpectPrepare(deleteDiscussionFlairAccessString)
	mock.ExpectPrepare(deleteDiscussionUserAccessString)
	mock.ExpectPrepare(getDiscussionInviteByIDString)
	mock.ExpectPrepare(getDiscussionRequestAccessByIDString)
	mock.ExpectPrepare(getDiscussionInvitesForUserString)
	mock.ExpectPrepare(getSentDiscussionInvitesForUserString)
	mock.ExpectPrepare(getDiscussionAccessRequestsString)
	mock.ExpectPrepare(getSentDiscussionAccessRequestsForUserString)
	mock.ExpectPrepare(getInviteLinksForDiscussion)
	mock.ExpectPrepare(putDiscussionInviteRecordString)
	mock.ExpectPrepare(putDiscussionAccessRequestString)
	mock.ExpectPrepare(updateDiscussionInviteRecordString)
	mock.ExpectPrepare(updateDiscussionAccessRequestString)
	mock.ExpectPrepare(upsertInviteLinksForDiscussion)
}

func mockPreparedStatementsWithError(mock sqlmock.Sqlmock) {
	mock.ExpectPrepare(getPostByIDString).WillReturnError(fmt.Errorf("error"))
}
