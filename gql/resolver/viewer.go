package resolver

import (
	"time"

	"github.com/graph-gophers/graphql-go"
)

type viewerResolver struct {
	id           graphql.ID
	discussionId graphql.ID
}

func (r *viewerResolver) ID() string {
	return "1"
}

func (r *viewerResolver) NotificationPreferences() (*discussionNotificationPreferencesResolver, error) {
	return &discussionNotificationPreferencesResolver{
		discussionId: r.discussionId,
	}, nil
}

func (r *viewerResolver) Discussion() (*discussionResolver, error) {
	return nil, nil
}

func (r *viewerResolver) LastViewed() *string {
	nowAsStr := time.Now().Format(time.RFC3339)
	return &nowAsStr
}

func (r *viewerResolver) LastPostViewed() (*postResolver, error) {
	return &postResolver{
		id: "1",
	}, nil
}

func (r *viewerResolver) Bookmarks() (*postsConnectionResolver, error) {
	return &postsConnectionResolver{
		ids: []graphql.ID{"1", "2"},
	}, nil
}

type viewersConnectionResolver struct {
	ids  []graphql.ID
	from graphql.ID
	to   graphql.ID
}

func (r *viewersConnectionResolver) TotalCount() int32 {
	return 0
}

func (r *viewersConnectionResolver) Edges() (*[]*viewersEdgeResolver, error) {
	return nil, nil
}

func (r *viewersConnectionResolver) PageInfo() (*pageInfoResolver, error) {
	return nil, nil
}

type viewersEdgeResolver struct {
	cursor graphql.ID
	id     graphql.ID
}

func (r *viewersEdgeResolver) Cursor() graphql.ID {
	return r.cursor
}

func (r *viewersEdgeResolver) Node() *viewerResolver {
	return &viewerResolver{
		id: r.id,
	}
}
