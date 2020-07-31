package backend

import (
	"context"
	"fmt"
	"strings"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/twitter"
	"github.com/delphis-inc/delphisbe/internal/util"
	"github.com/sirupsen/logrus"
)

const (
	/* We assume that the most interesting autocompletes are in the first page.
	   A more sophisticated connection-based fetch would be better, but would
	   also add unnecessary overhead. */
	twitterAutocompletesPageSize = 20
	twitterAutocompletsMaxPages  = 1
)

// We use second user's token here for rate limiting reasons.
func (d *delphisBackend) DoesTwitterUserFollowUser(ctx context.Context, twitterClient twitter.TwitterClient, firstUser model.SocialInfo, secondUser model.SocialInfo) (bool, error) {
	if secondUser.Network != util.SocialNetworkTwitter || firstUser.Network != util.SocialNetworkTwitter {
		return false, fmt.Errorf("Both users must be twitter accounts")
	}

	relationship, err := twitterClient.FriendshipLookup(secondUser.ScreenName, firstUser.ScreenName)
	if err != nil || relationship == nil {
		return false, fmt.Errorf("Failed contacting Twitter")
	}

	return relationship.Target.Following, nil
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
		if strings.ToLower(info.Network) == util.SocialNetworkTwitter {
			accessToken = info.AccessToken
			accessTokenSecret = info.AccessTokenSecret
		}
	}

	return accessToken, accessTokenSecret, nil
}

func (d *delphisBackend) GetTwitterClientWithUserTokens(ctx context.Context) (twitter.TwitterClient, error) {
	/* Obtain infos needed for creating Twitter API client */
	accessToken, accessTokenSecret, err := d.GetTwitterAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	return d.GetTwitterClientWithAccessTokens(ctx, accessToken, accessTokenSecret)
}

func (d *delphisBackend) GetTwitterClientWithAccessTokens(ctx context.Context, accessToken string, accessTokenSecret string) (twitter.TwitterClient, error) {
	consumerKey := d.config.Twitter.ConsumerKey
	consumerSecret := d.config.Twitter.ConsumerSecret
	/* Check that everything is ready to go */
	return d.twitterBackend.GetTwitterClientWithAccessTokens(ctx, consumerKey, consumerSecret, accessToken, accessTokenSecret)
}

func (d *delphisBackend) GetTwitterUserHandleAutocompletes(ctx context.Context, twitterClient twitter.TwitterClient, query string, discussionID string, invitingParticipantID string) ([]*model.TwitterUserInfo, error) {
	/* Check the list of all the Twitter users already invited in the discussion for this participant.
	   Fetching all of them in a single query enhances the scalability of this function. */
	invitedTwitterHandles, err := d.db.GetInvitedTwitterHandlesByDiscussionIDAndInviterID(ctx, discussionID, invitingParticipantID)
	if err != nil {
		return nil, err
	}

	/* Fetch autocompletes result eagerly from twitter APIs. A connection-based paging
	   system would have more quality but would also introduce additional overhead.
	   As a tradeoff we limit the number of pages fetched by assuming that the best
	   results will be on top of the list */
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
				DisplayName:     twitterUser.Name,
				Name:            twitterUser.ScreenName,
				Verified:        twitterUser.Verified,
				ProfileImageURL: twitterUser.ProfileImageURLHttps,
				Invited:         isInvited,
			}
			results = append(results, twitterUserInfo)
		}
	}

	return results, nil
}

func (d *delphisBackend) InviteTwitterUsersToDiscussion(ctx context.Context, twitterClient twitter.TwitterClient, twitterUserInfos []*model.TwitterUserInput, discussionID, invitingParticipantID string) ([]*model.DiscussionInvite, error) {
	var invitations []*model.DiscussionInvite

	/* Check that the user is autenticated */
	authedUser := auth.GetAuthedUser(ctx)
	if authedUser == nil {
		return nil, fmt.Errorf("Need auth")
	}

	/* Obtain all the Twitter users that this participant has already invited in the discussion */
	alreadyInvitedHandles, err := d.db.GetInvitedTwitterHandlesByDiscussionIDAndInviterID(ctx, discussionID, invitingParticipantID)
	if err != nil {
		return nil, err
	}

	/* Filter only the Twitter users that haven't already been invited */
	var screenNames []string
	lenAlreadyInvitedHandles := len(alreadyInvitedHandles)
	lenScreenNames := 0
	for _, u := range twitterUserInfos {
		alreadyInvited := false
		for i := 0; i < lenAlreadyInvitedHandles && !alreadyInvited; i++ {
			if *alreadyInvitedHandles[i] == u.Name {
				alreadyInvited = true
			}
		}
		if !alreadyInvited {
			screenNames = append(screenNames, u.Name)
			lenScreenNames++
		}
	}

	/* Avoid useless logic */
	if lenScreenNames == 0 {
		return invitations, nil
	}

	/* Leverage Twitter APIs Lookup query to retrieve users in batch with a single request */
	twitterUsers, err := twitterClient.LookupUsers(screenNames)
	if err != nil {
		return nil, err
	}

	/* Iterate throug twitter users and send them individual invitations */
	for _, twitterUser := range twitterUsers {

		/* Check that the user is effectively one of the desired ones */
		isCorrectUser := false
		for i := 0; i < lenScreenNames && !isCorrectUser; i++ {
			if screenNames[i] == twitterUser.ScreenName {
				isCorrectUser = true
			}
		}

		if isCorrectUser {
			/* Get invited user. If the user is not present in the system, we create it
			with a dummy access token. Note, the datastore will not overwrite the tokens
			with the dummy ones if valid tokens are already present */
			userObj, err := d.GetOrCreateUser(ctx, LoginWithTwitterInput{
				User:              &twitterUser,
				AccessToken:       "",
				AccessTokenSecret: "",
			}, nil)
			if err != nil {
				logrus.WithError(err).Errorf("Got an error creating a user")
				return nil, err
			}

			/* Prevent users from inviting themselves */
			if userObj.ID != authedUser.UserID {
				invite, err := d.InviteUserToDiscussion(ctx, userObj.ID, discussionID, invitingParticipantID)
				if err != nil {
					return nil, err
				}
				invitations = append(invitations, invite)
				/* TODO: (?) We may consider to notify users in some way external to the app, like email (if public) or twitter
				   dm (if they follow the authed user), in order to invite users to install the app.
				   This would be the place to do it. */
			}
		}
	}

	return invitations, nil
}
