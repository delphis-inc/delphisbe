package backend

import (
	"context"
	"time"

	"github.com/nedrocks/delphisbe/graph/model"
	"github.com/nedrocks/delphisbe/internal/util"
)

func (d *delphisBackend) CreatePost(ctx context.Context, discussionKey model.DiscussionParticipantKey, content string) (*model.Post, error) {
	postContent := model.PostContent{
		ID:      util.UUIDv4(),
		Content: content,
	}

	post := model.Post{
		ID:            util.UUIDv4(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		DiscussionID:  discussionKey.DiscussionID,
		ParticipantID: discussionKey.ParticipantID,
		PostContentID: postContent.ID,
		PostContent:   postContent,
	}

	postObj, err := d.db.PutPost(ctx, post)

	if err != nil {
		return nil, err
	}

	return postObj, nil
}

func (d *delphisBackend) GetPostsByDiscussionID(ctx context.Context, discussionID string) ([]*model.Post, error) {
	return d.db.GetPostsByDiscussionID(ctx, discussionID)
}
