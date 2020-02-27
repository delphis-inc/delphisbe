package resolver

import "github.com/graph-gophers/graphql-go"

type Resolver struct{}

func (r *Resolver) Discussion(args struct{ Id graphql.ID }) *discussionResolver {
	return &discussionResolver{
		id: args.Id,
	}
}

func (r *Resolver) Post(args struct{ Id graphql.ID }) *postResolver {
	return &postResolver{
		id: args.Id,
	}
}

func (r *Resolver) User(args struct{ Id graphql.ID }) *userResolver {
	return &userResolver{
		id: args.Id,
	}
}

func (r *Resolver) Viewer(args struct{ Id graphql.ID }) *viewerResolver {
	return &viewerResolver{
		id: args.Id,
	}
}
