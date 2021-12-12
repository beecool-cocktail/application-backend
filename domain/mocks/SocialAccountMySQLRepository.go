// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/beecool-cocktail/application-backend/domain"
	mock "github.com/stretchr/testify/mock"
)

// SocialAccountMySQLRepository is an autogenerated mock type for the SocialAccountMySQLRepository type
type SocialAccountMySQLRepository struct {
	mock.Mock
}

// QueryById provides a mock function with given fields: ctx, id
func (_m *SocialAccountMySQLRepository) QueryById(ctx context.Context, id string) (*domain.SocialAccount, error) {
	ret := _m.Called(ctx, id)

	var r0 *domain.SocialAccount
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.SocialAccount); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.SocialAccount)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, s, u
func (_m *SocialAccountMySQLRepository) Store(ctx context.Context, s *domain.SocialAccount, u *domain.User) error {
	ret := _m.Called(ctx, s, u)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.SocialAccount, *domain.User) error); ok {
		r0 = rf(ctx, s, u)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
