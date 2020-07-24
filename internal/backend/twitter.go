package backend

import (
	"context"
	"fmt"
	"strings"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/sirupsen/logrus"
)

const (
	/* We assume that the most interesting autocompletes are in the first page.
	   A more sophisticated connection-based fetch would be better, but would
	   also add unnecessary overhead. */
	twitterAutocompletesPageSize = 20
	twitterAutocompletsMaxPages  = 1
)

/* This is an interface that abstracts the business logic related to Twitter APIs.
   Having an internal interface helps reducing dependency binding and helps in testing. */
type TwitterClient interface {
	SearchUsers(query string, page int, count int) ([]twitter.User, error)
	LookupUsers(screenNames []string) ([]twitter.User, error)
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

func (d *delphisBackend) GetTwitterAccessToken(ctx context.Context) (string, string, error) {
	/* Get the authed user */
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return "", "", fmt.Errorf("Need auth")
	}

	/* Obtain authed user profile */
	authedUserProfile, err := d.db.GetUserProfileByUserID(ctx, authedUser.UserID)
	if err != nil {
		return "", "", err
	}

	/* Obtain authed user social info  */
	authedSocialInfo, err := d.db.GetSocialInfosByUserProfileID(ctx, *&authedUserProfile.ID)
	if err != nil {
		return "", "", err
	}

	/* Extract tokens from social info */
	accessToken := ""
	accessTokenSecret := ""
	for _, info := range authedSocialInfo {
		if strings.ToLower(info.Network) == "twitter" {
			accessToken = info.AccessToken
			accessTokenSecret = info.AccessTokenSecret
		}
	}

	return accessToken, accessTokenSecret, nil
}

func (d *delphisBackend) GetTwitterClientWithUserTokens(ctx context.Context) (TwitterClient, error) {
	/* Obtain infos needed for creating Twitter API client */
	consumerKey := d.config.Twitter.ConsumerKey
	consumerSecret := d.config.Twitter.ConsumerSecret
	accessToken, accessTokenSecret, err := d.GetTwitterAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	/* Check that everything is ready to go */
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

func (d *delphisBackend) GetTwitterUserHandleAutocompletes(ctx context.Context, twitterClient TwitterClient, query string, discussionID string, invitingParticipantID string) ([]*model.TwitterUserInfo, error) {
	/* Fetch autocompletes result eagerly from twitter APIs. A connection-based paging
	   system would have more quality but would also introduce additional overhead.
	   As a tradeoff we limit the number of pages fetched by assuming that the best
	   results will be on top of the list */
	invitedTwitterHandles, err := d.db.GetInvitedTwitterHandlesByDiscussionIDAndInviterID(ctx, discussionID, invitingParticipantID)
	if err != nil {
		return nil, err
	}

	lenInvitedTwitterHandles := len(invitedTwitterHandles)
	var results []*model.TwitterUserInfo
	curPage := 0
	for resultSize := 0; (curPage == 0 || resultSize == twitterAutocompletesPageSize) && curPage < twitterAutocompletsMaxPages; curPage++ {
		twitterUsers, err := twitterClient.SearchUsers(query, curPage, twitterAutocompletesPageSize)
		if err != nil {
			return nil, err
		}
		resultSize = len(twitterUsers)
		for _, twitterUser := range twitterUsers {
			isInvited := false
			for i := 0; i < lenInvitedTwitterHandles && !isInvited; i++ {
				if *invitedTwitterHandles[i] == twitterUser.ScreenName {
					isInvited = true
				}
			}
			twitterUserInfo := &model.TwitterUserInfo{
				ID:              twitterUser.IDStr,
				DiplayName:      twitterUser.Name,
				Name:            twitterUser.ScreenName,
				IsVerified:      twitterUser.Verified,
				ProfileImageURL: twitterUser.ProfileImageURLHttps,
				IsInvited:       isInvited,
			}
			results = append(results, twitterUserInfo)
		}
	}

	return results, nil
}

func (d *delphisBackend) InviteTwitterUsersToDiscussion(ctx context.Context, twitterClient TwitterClient, twitterUserInfos []*model.TwitterUserInput, discussionID, invitingParticipantID string) ([]*model.DiscussionInvite, error) {
	/* Check that the user is autenticated */
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	var screenNames []string
	for _, u := range twitterUserInfos {
		screenNames = append(screenNames, u.Name)
	}
	/* Leverage Twitter APIs Lookup query to retrieve users in batch with a single request */
	twitterUsers, err := twitterClient.LookupUsers(screenNames)
	if err != nil {
		return nil, err
	}

	/* Iterate throug twitter users and send them individual invitations */
	var invitations []*model.DiscussionInvite
	for _, twitterUser := range twitterUsers {
		/* Get invited user. If the user is not present in the system, we create it
		with a dummy access token. Note, the datastore will not overwrite the tokens
		with the dummy ones if valid tokens are already present */
		userObj, err := d.GetOrCreateUser(ctx, LoginWithTwitterInput{
			User:              &twitterUser,
			AccessToken:       "",
			AccessTokenSecret: "",
		})
		if err != nil {
			logrus.WithError(err).Errorf("Got an error creating a user")
			return nil, err
		}

		/* Prevent users from inviting themselves */
		if userObj.ID != authedUser.UserID {
			/* Verify that an invite is not already present for such an user
			   NOTE: Should we check for already accepted invitations too? Maybe we can check if the user
			   is already a participant even before calling this function. */
			userInvites, err := d.GetDiscussionInvitesByUserIDAndStatus(ctx, userObj.ID, model.InviteRequestStatusPending)
			if err != nil {
				return nil, err
			}

			/* If the user has already a pending invitation, we return it instead of creating a new one */
			if len(userInvites) == 0 {
				invite, err := d.InviteUserToDiscussion(ctx, userObj.ID, discussionID, invitingParticipantID)
				if err != nil {
					return nil, err
				}
				invitations = append(invitations, invite)
				/* TODO: (?) We may consider to notify users in some way external to the app, like email (if public) or twitter
				   dm (if they follow the authed user), in order to invite users to install the app. */
			} else {
				invitations = append(invitations, userInvites[0])
			}
		}
	}

	return invitations, nil
}
