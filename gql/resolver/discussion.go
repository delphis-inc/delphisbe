package resolver

import (
	"github.com/graph-gophers/graphql-go"
)

type discussionResolver struct {
	id            graphql.ID
	anonymityType string
}

func (r *discussionResolver) ID() graphql.ID {
	return r.id
}

func (r *discussionResolver) AnonymityType() string {
	return r.anonymityType
}

type discussionsConnectionResolver struct {
	ids  []graphql.ID
	from graphql.ID
	to   graphql.ID
}

func (r *discussionsConnectionResolver) TotalCount() int32 {
	return 0
}

func (r *discussionsConnectionResolver) Edges() (*[]*discussionsEdgeResolver, error) {
	return nil, nil
}

func (r *discussionsConnectionResolver) PageInfo() (*pageInfoResolver, error) {
	return nil, nil
}

type discussionsEdgeResolver struct {
	cursor graphql.ID
	id     graphql.ID
}

func (r *discussionsEdgeResolver) Cursor() graphql.ID {
	return r.cursor
}

func (r *discussionsEdgeResolver) Node() *discussionResolver {
	return &discussionResolver{
		id: r.id,
	}
}
