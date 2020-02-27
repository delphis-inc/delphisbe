package resolver

import (
	"time"

	"github.com/graph-gophers/graphql-go"
)

type postBookmarkResolver struct {
	id graphql.ID
}

func (r *postBookmarkResolver) ID() graphql.ID {
	return r.id
}

func (r *postBookmarkResolver) Discussion() (*discussionResolver, error) {
	return &discussionResolver{
		id: "1",
	}, nil
}

func (r *postBookmarkResolver) Post() (*postResolver, error) {
	return &postResolver{
		id: "1",
	}, nil
}

func (r *postBookmarkResolver) CreatedAt() string {
	return time.Now().Format(time.RFC3339)
}

type postBookmarksConnectionResolver struct {
	ids  []graphql.ID
	from graphql.ID
	to   graphql.ID
}

func (r *postBookmarksConnectionResolver) TotalCount() int32 {
	return 0
}

func (r *postBookmarksConnectionResolver) Edges() (*[]*postBookmarksEdgeResolver, error) {
	return nil, nil
}

func (r *postBookmarksConnectionResolver) PageInfo() (*pageInfoResolver, error) {
	return nil, nil
}

type postBookmarksEdgeResolver struct {
	cursor graphql.ID
	id     graphql.ID
}

func (r *postBookmarksEdgeResolver) Cursor() graphql.ID {
	return r.cursor
}

func (r *postBookmarksEdgeResolver) Node() *postBookmarkResolver {
	return &postBookmarkResolver{
		id: r.id,
	}
}
