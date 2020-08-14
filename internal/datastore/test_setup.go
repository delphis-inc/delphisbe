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
	mock.ExpectPrepare(getDiscussionByLinkSlugString)
	mock.ExpectPrepare(getDiscussionArchiveByDiscussionIDString)
	mock.ExpectPrepare(upsertDiscussionArchiveString)
	mock.ExpectPrepare(getModeratorByUserIDString)
	mock.ExpectPrepare(getModeratorByUserIDAndDiscussionIDString)
	mock.ExpectPrepare(getModeratedDiscussionsByUserIDString)
	mock.ExpectPrepare(getDiscussionsByUserAccessString)
	mock.ExpectPrepare(getDiscussionUserAccessString)
	mock.ExpectPrepare(getDUAForEverythingNotificationsString)
	mock.ExpectPrepare(getDUAForMentionNotificationsString)
	mock.ExpectPrepare(upsertDiscussionUserAccessString)
	mock.ExpectPrepare(deleteDiscussionUserAccessString)
	mock.ExpectPrepare(getDiscussionRequestAccessByIDString)
	mock.ExpectPrepare(getDiscussionAccessRequestsString)
	mock.ExpectPrepare(getDiscussionAccessRequestByUserIDString)
	mock.ExpectPrepare(getSentDiscussionAccessRequestsForUserString)
	mock.ExpectPrepare(putDiscussionAccessRequestString)
	mock.ExpectPrepare(updateDiscussionAccessRequestString)
	mock.ExpectPrepare(getAccessLinkBySlugString)
	mock.ExpectPrepare(getAccessLinkByDiscussionIDString)
	mock.ExpectPrepare(putAccessLinkForDiscussionString)
	mock.ExpectPrepare(getNextShuffleTimeForDiscussionIDString)
	mock.ExpectPrepare(putNextShuffleTimeForDiscussionIDString)
	mock.ExpectPrepare(getDiscussionsToShuffle)
	mock.ExpectPrepare(incrDiscussionShuffleCount)
	mock.ExpectPrepare(getViewerForDiscussionIDUserID)
	mock.ExpectPrepare(updateViewerLastViewed)
}

func mockPreparedStatementsWithError(mock sqlmock.Sqlmock) {
	mock.ExpectPrepare(getPostByIDString).WillReturnError(fmt.Errorf("error"))
}
