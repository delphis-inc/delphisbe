package backend

import (
	"context"
	"fmt"
	"strings"

	"github.com/delphis-inc/delphisbe/graph/model"
	"github.com/delphis-inc/delphisbe/internal/auth"
	"github.com/delphis-inc/delphisbe/internal/twitter"
	"github.com/delphis-inc/delphisbe/internal/util"
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
	/* Fetch autocompletes result eagerly from twitter APIs. A connection-based paging
	   system would have more quality but would also introduce additional overhead.
	   As a tradeoff we limit the number of pages fetched by assuming that the best
	   results will be on top of the list */
	var results []*model.TwitterUserInfo
	curPage := 0
	for resultSize := 0; (curPage == 0 || resultSize == twitterAutocompletesPageSize) && curPage < twitterAutocompletsMaxPages; curPage++ {
		twitterUsers, err := twitterClient.SearchUsers(query, curPage, twitterAutocompletesPageSize)
		if err != nil {
			return nil, err
		}
		resultSize = len(twitterUsers)
		for _, twitterUser := range twitterUsers {
			twitterUserInfo := &model.TwitterUserInfo{
				ID:              twitterUser.IDStr,
				DisplayName:     twitterUser.Name,
				Name:            twitterUser.ScreenName,
				Verified:        twitterUser.Verified,
				ProfileImageURL: twitterUser.ProfileImageURLHttps,
				Invited:         false,
			}
			results = append(results, twitterUserInfo)
		}
	}

	return results, nil
}
