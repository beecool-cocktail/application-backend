// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/beecool-cocktail/application-backend/domain"
	mock "github.com/stretchr/testify/mock"
)

// UserUsecase is an autogenerated mock type for the UserUsecase type
type UserUsecase struct {
	mock.Mock
}

// Logout provides a mock function with given fields: ctx, id
func (_m *UserUsecase) Logout(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueryById provides a mock function with given fields: ctx, id
func (_m *UserUsecase) QueryById(ctx context.Context, id int64) (domain.User, error) {
	ret := _m.Called(ctx, id)

	var r0 domain.User
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.User); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUserInfo provides a mock function with given fields: ctx, d, ui
func (_m *UserUsecase) UpdateUserInfo(ctx context.Context, d *domain.User, ui *domain.UserImage) error {
	ret := _m.Called(ctx, d, ui)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.User, *domain.UserImage) error); ok {
		r0 = rf(ctx, d, ui)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
