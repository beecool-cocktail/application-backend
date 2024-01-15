// Code generated by mockery v2.30.16. DO NOT EDIT.

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

// QueryUserNameByID provides a mock function with given fields: ctx, id
func (_m *UserRedisRepository) QueryUserNameByID(ctx context.Context, id int64) (string, error) {
	ret := _m.Called(ctx, id)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (string, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) string); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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

// NewUserRedisRepository creates a new instance of UserRedisRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRedisRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRedisRepository {
	mock := &UserRedisRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
