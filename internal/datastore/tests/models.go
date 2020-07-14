package tests

import (
	"github.com/delphis-inc/delphisbe/graph/model"
)

//const (
//	Discussion1ID = "77aceb56-5dcd-4ac8-8833-bbfd058aa882"
//	Discussion2ID = "bf195c5d-1f02-478a-9a44-7dedaca4945b"
//	Participant1ID        = "d707df47-be14-4890-ba03-3276cc5388b1"
//	Participant2ID        = "4d3ab11b-76b7-468b-9f6d-7d0deef20ced"
//)

var (
	Discussion1ID  = "77aceb56-5dcd-4ac8-8833-bbfd058aa882"
	Discussion2ID  = "bf195c5d-1f02-478a-9a44-7dedaca4945b"
	Participant1ID = "d707df47-be14-4890-ba03-3276cc5388b1"
	Participant2ID = "4d3ab11b-76b7-468b-9f6d-7d0deef20ced"
	Post1ID        = "b39295be-1cf6-4777-a558-6166494cdc7c"
	Post2ID        = "288e13cb-ace6-496a-ae15-7f69bbbbe156"
	Post3ID        = "85272a83-7814-4752-887b-6b700add1d02"
	PostContent1ID = "882a2798-1a6b-4a94-9939-dbea3ede19d5"
	PostContent2ID = "a479703d-442b-421b-81a5-e04fdaf43e73"
	PostContent3ID = "ae844838-b1bd-4cfb-a5f8-bd7b46e89a33"
	testViewer     = "testViewer"
	testUser       = "testUser"
)

func ValidDiscussions() []model.Discussion {
	return []model.Discussion{
		{
			ID:            Discussion1ID,
			Title:         "Discussion1Test",
			AnonymityType: "WEAK",
		},
		{
			ID:            Discussion2ID,
			Title:         "Discussion2Test",
			AnonymityType: "WEAK",
		},
	}
}

func ValidParticipants() []model.Participant {
	return []model.Participant{
		{
			ID:            Participant1ID,
			ParticipantID: 0,
			DiscussionID:  &Discussion1ID,
			ViewerID:      &testViewer,
			UserID:        &testUser,
		},
		{
			ID:            Participant2ID,
			ParticipantID: 0,
			DiscussionID:  &Discussion1ID,
			ViewerID:      &testViewer,
			UserID:        &testUser,
		},
	}
}

func ValidPosts() []model.Post {
	return []model.Post{
		{
			ID:            Post1ID,
			DiscussionID:  &Discussion1ID,
			ParticipantID: &Participant1ID,
			PostContent: &model.PostContent{
				ID:      PostContent1ID,
				Content: "hi, from post 1",
			},
		},
		{
			ID:            Post2ID,
			DiscussionID:  &Discussion1ID,
			ParticipantID: &Participant2ID,
			PostContent: &model.PostContent{
				ID:      PostContent2ID,
				Content: "hello, post 1",
			},
		},
		{
			ID:            Post3ID,
			DiscussionID:  &Discussion1ID,
			ParticipantID: &Participant1ID,
			PostContent: &model.PostContent{
				ID:      PostContent3ID,
				Content: "nic to meet you participant 2",
			},
			QuotedPostID: &Post2ID,
		},
	}
}
