package resolver

import "github.com/graph-gophers/graphql-go"

type postResolver struct {
	id graphql.ID
}

type postsConnectionResolver struct {
	ids  []graphql.ID
	from graphql.ID
	to   graphql.ID
}

func (r *postsConnectionResolver) TotalCount() int32 {
	return 0
}

func (r *postsConnectionResolver) Edges() (*[]*postsEdgeResolver, error) {
	return nil, nil
}

func (r *postsConnectionResolver) PageInfo() (*pageInfoResolver, error) {
	return nil, nil
}

type postsEdgeResolver struct {
	cursor graphql.ID
	id     graphql.ID
}

func (r *postsEdgeResolver) Cursor() graphql.ID {
	return r.cursor
}

func (r *postsEdgeResolver) Node() *postResolver {
	return &postResolver{
		id: r.id,
	}
}
