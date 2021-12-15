// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/beecool-cocktail/application-backend/domain"
	mock "github.com/stretchr/testify/mock"

	oauth2 "golang.org/x/oauth2"
)

// SocialAccountGoogleOAuthRepository is an autogenerated mock type for the SocialAccountGoogleOAuthRepository type
type SocialAccountGoogleOAuthRepository struct {
	mock.Mock
}

// Exchange provides a mock function with given fields: ctx, code
func (_m *SocialAccountGoogleOAuthRepository) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	ret := _m.Called(ctx, code)

	var r0 *oauth2.Token
	if rf, ok := ret.Get(0).(func(context.Context, string) *oauth2.Token); ok {
		r0 = rf(ctx, code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*oauth2.Token)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserInfo provides a mock function with given fields: ctx, token
func (_m *SocialAccountGoogleOAuthRepository) GetUserInfo(ctx context.Context, token *oauth2.Token) (*domain.GoogleUserInfo, error) {
	ret := _m.Called(ctx, token)

	var r0 *domain.GoogleUserInfo
	if rf, ok := ret.Get(0).(func(context.Context, *oauth2.Token) *domain.GoogleUserInfo); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.GoogleUserInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *oauth2.Token) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}