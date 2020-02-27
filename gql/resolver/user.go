package resolver

import "github.com/graph-gophers/graphql-go"

type userResolver struct {
	id graphql.ID
}

func (r *userResolver) ID() graphql.ID {
	return r.id
}

func (r *userResolver) Participants() (*participantsConnectionResolver, error) {
	return &participantsConnectionResolver{
		ids: []graphql.ID{"1", "2"},
	}, nil
}

func (r *userResolver) Viewers() (*viewersConnectionResolver, error) {
	return &viewersConnectionResolver{
		ids: []graphql.ID{"1", "2"},
	}, nil
}
