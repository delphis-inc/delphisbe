package datastore

import (
	"github.com/nedrocks/delphisbe/graph/model"
)

type postTestData struct {
	discussionID string
	post         model.Post
}

//func TestPostDatastore(t *testing.T) {
//	// Testing variables
//	data := TestData{
//		Discussions:  tests.ValidDiscussions(),
//		Participants: tests.ValidParticipants(),
//		Posts:        tests.ValidPosts(),
//	}
//	postObj := model.Post{
//		ID:            "3031f9ee-fa75-4004-8572-0eba5a0a43b7",
//		DiscussionID:  &tests.Discussion1ID,
//		ParticipantID: &tests.Participant2ID,
//		PostContent: &model.PostContent{
//			ID:      "6b67e749-1edd-4f0f-b8e6-732d28ce3da9",
//			Content: "unit test post",
//		},
//	}
//
//	// Initialize the test DB once per file
//	db, close, err := MakeDatastore(context.Background(), data)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer close()
//
//	// Loop over table tests
//	tests := []struct {
//		scenario string
//		data     postTestData
//		test     func(ctx context.Context, t *testing.T, db Datastore, data postTestData)
//	}{
//		{
//			scenario: "insert new post",
//			data:     postTestData{post: postObj},
//			test: func(ctx context.Context, t *testing.T, db Datastore, data postTestData) {
//				tx, _ := db.BeginTx(ctx)
//
//				if err := db.PutPostContent(ctx, tx, *data.post.PostContent); err != nil {
//					t.Fatal(err)
//				}
//				resp, err := db.PutPost(ctx, tx, data.post)
//				assert.NoError(t, err)
//				assert.Equal(t, postObj.ID, resp.ID)
//
//				// Keep db clean
//				db.RollbackTx(ctx, tx)
//			},
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.scenario, func(t *testing.T) {
//			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
//			defer cancel()
//
//			test.test(ctx, t, db, test.data)
//		})
//	}
//}
