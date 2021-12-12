// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/beecool-cocktail/application-backend/domain"
	mock "github.com/stretchr/testify/mock"
)

// UserMySQLRepository is an autogenerated mock type for the UserMySQLRepository type
type UserMySQLRepository struct {
	mock.Mock
}

// QueryById provides a mock function with given fields: ctx, id
func (_m *UserMySQLRepository) QueryById(ctx context.Context, id int64) (*domain.User, error) {
	ret := _m.Called(ctx, id)

	var r0 *domain.User
	if rf, ok := ret.Get(0).(func(context.Context, int64) *domain.User); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, d
func (_m *UserMySQLRepository) Store(ctx context.Context, d *domain.User) error {
	ret := _m.Called(ctx, d)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.User) error); ok {
		r0 = rf(ctx, d)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
