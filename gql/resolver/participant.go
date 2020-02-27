package resolver

import (
	"github.com/graph-gophers/graphql-go"
)

type participantResolver struct {
	id graphql.ID
}

func (r *participantResolver) ID() graphql.ID {
	return r.id
}

func (r *participantResolver) Discussion() (*discussionResolver, error) {
	return &discussionResolver{
		id: "1",
	}, nil
}

func (r *participantResolver) Viewer() (*viewerResolver, error) {
	return &viewerResolver{
		id: "1",
	}, nil
}

func (r *participantResolver) DiscussionNotificationPreferences() (*notificationPreferencesResolver, error) {
	return &notificationPreferencesResolver{}, nil
}

func (r *participantResolver) Posts() (*postsConnectionResolver, error) {
	return &postsConnectionResolver{
		ids: []graphql.ID{"1", "2"},
	}, nil
}

type participantsConnectionResolver struct {
	ids  []graphql.ID
	from graphql.ID
	to   graphql.ID
}

func (r *participantsConnectionResolver) TotalCount() int32 {
	return 0
}

func (r *participantsConnectionResolver) Edges() (*[]*participantsEdgeResolver, error) {
	return nil, nil
}

func (r *participantsConnectionResolver) PageInfo() (*pageInfoResolver, error) {
	return nil, nil
}

type participantsEdgeResolver struct {
	cursor graphql.ID
	id     graphql.ID
}

func (r *participantsEdgeResolver) Cursor() graphql.ID {
	return r.cursor
}

func (r *participantsEdgeResolver) Node() *participantResolver {
	return &participantResolver{
		id: r.id,
	}
}
