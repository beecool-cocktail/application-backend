// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/beecool-cocktail/application-backend/domain"
	mock "github.com/stretchr/testify/mock"
)

// UserRedisRepository is an autogenerated mock type for the UserRedisRepository type
type UserRedisRepository struct {
	mock.Mock
}

// Store provides a mock function with given fields: ctx, r
func (_m *UserRedisRepository) Store(ctx context.Context, r *domain.UserCache) error {
	ret := _m.Called(ctx, r)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.UserCache) error); ok {
		r0 = rf(ctx, r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateBasicInfo provides a mock function with given fields: ctx, r
func (_m *UserRedisRepository) UpdateBasicInfo(ctx context.Context, r *domain.UserCache) error {
	ret := _m.Called(ctx, r)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.UserCache) error); ok {
		r0 = rf(ctx, r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateToken provides a mock function with given fields: ctx, r
func (_m *UserRedisRepository) UpdateToken(ctx context.Context, r *domain.UserCache) error {
	ret := _m.Called(ctx, r)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.UserCache) error); ok {
		r0 = rf(ctx, r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
