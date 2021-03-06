// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	twitter "github.com/delphis-inc/delphisbe/internal/twitter"
	mock "github.com/stretchr/testify/mock"
)

// TwitterBackend is an autogenerated mock type for the TwitterBackend type
type TwitterBackend struct {
	mock.Mock
}

// GetTwitterClientWithAccessTokens provides a mock function with given fields: ctx, consumerKey, consumerSecret, accessToken, accessTokenSecret
func (_m *TwitterBackend) GetTwitterClientWithAccessTokens(ctx context.Context, consumerKey string, consumerSecret string, accessToken string, accessTokenSecret string) (twitter.TwitterClient, error) {
	ret := _m.Called(ctx, consumerKey, consumerSecret, accessToken, accessTokenSecret)

	var r0 twitter.TwitterClient
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) twitter.TwitterClient); ok {
		r0 = rf(ctx, consumerKey, consumerSecret, accessToken, accessTokenSecret)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(twitter.TwitterClient)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string) error); ok {
		r1 = rf(ctx, consumerKey, consumerSecret, accessToken, accessTokenSecret)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
