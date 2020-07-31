package twitter

import (
	"context"
	"fmt"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

/* This is an interface that abstracts the business logic related to Twitter APIs.
   Having an internal interface helps reducing dependency binding and helps in testing. */
type TwitterClient interface {
	SearchUsers(query string, page int, count int) ([]twitter.User, error)
	LookupUsers(screenNames []string) ([]twitter.User, error)
	FriendshipLookup(fromScreenName, toScreenName string) (*twitter.Relationship, error)
}

/* Implementation of the interface above based on an external package */
type twitterClient struct {
	client *twitter.Client
}

func (t *twitterClient) SearchUsers(query string, page int, count int) ([]twitter.User, error) {
	userSearchParams := &twitter.UserSearchParams{
		Query: query,
		Page:  page,
		Count: count,
	}
	twitterUsers, _, err := t.client.Users.Search(query, userSearchParams)
	return twitterUsers, err
}

func (t *twitterClient) LookupUsers(screenNames []string) ([]twitter.User, error) {
	twitterUsers, _, err := t.client.Users.Lookup(&twitter.UserLookupParams{
		ScreenName: screenNames,
	})
	return twitterUsers, err
}

func (t *twitterClient) FriendshipLookup(fromScreenName, toScreenName string) (*twitter.Relationship, error) {
	relationship, _, err := t.client.Friendships.Show(&twitter.FriendshipShowParams{
		SourceScreenName: fromScreenName,
		TargetScreenName: toScreenName,
	})

	return relationship, err
}

type TwitterBackend interface {
	GetTwitterClientWithAccessTokens(ctx context.Context, consumerKey string, consumerSecret string, accessToken string, accessTokenSecret string) (TwitterClient, error)
}

type TwitterBackendImpl struct{}

func (t *TwitterBackendImpl) GetTwitterClientWithAccessTokens(ctx context.Context, consumerKey string, consumerSecret string, accessToken string, accessTokenSecret string) (TwitterClient, error) {
	if len(consumerKey) == 0 || len(consumerSecret) == 0 || len(accessToken) == 0 || len(accessTokenSecret) == 0 {
		return nil, fmt.Errorf("There is a problem retrieving authed user Twitter data")
	}

	/* Create client object */
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return &twitterClient{
		client: client,
	}, nil
}
