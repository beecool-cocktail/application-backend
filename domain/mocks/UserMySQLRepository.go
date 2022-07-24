// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/beecool-cocktail/application-backend/domain"
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

// UserMySQLRepository is an autogenerated mock type for the UserMySQLRepository type
type UserMySQLRepository struct {
	mock.Mock
}

// QueryById provides a mock function with given fields: ctx, id
func (_m *UserMySQLRepository) QueryById(ctx context.Context, id int64) (domain.User, error) {
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

// UpdateBasicInfo provides a mock function with given fields: ctx, d
func (_m *UserMySQLRepository) UpdateBasicInfo(ctx context.Context, d *domain.User) (int64, error) {
	ret := _m.Called(ctx, d)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *domain.User) int64); ok {
		r0 = rf(ctx, d)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.User) error); ok {
		r1 = rf(ctx, d)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateBasicInfoTx provides a mock function with given fields: ctx, tx, d
func (_m *UserMySQLRepository) UpdateBasicInfoTx(ctx context.Context, tx *gorm.DB, d *domain.User) (int64, error) {
	ret := _m.Called(ctx, tx, d)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.User) int64); ok {
		r0 = rf(ctx, tx, d)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *domain.User) error); ok {
		r1 = rf(ctx, tx, d)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateImage provides a mock function with given fields: ctx, d
func (_m *UserMySQLRepository) UpdateImage(ctx context.Context, d *domain.UserImage) (int64, error) {
	ret := _m.Called(ctx, d)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *domain.UserImage) int64); ok {
		r0 = rf(ctx, d)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *domain.UserImage) error); ok {
		r1 = rf(ctx, d)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateImageTx provides a mock function with given fields: ctx, tx, d
func (_m *UserMySQLRepository) UpdateImageTx(ctx context.Context, tx *gorm.DB, d *domain.UserImage) (int64, error) {
	ret := _m.Called(ctx, tx, d)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.UserImage) int64); ok {
		r0 = rf(ctx, tx, d)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *domain.UserImage) error); ok {
		r1 = rf(ctx, tx, d)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateNumberOfDraftTx provides a mock function with given fields: ctx, tx, d
func (_m *UserMySQLRepository) UpdateNumberOfDraftTx(ctx context.Context, tx *gorm.DB, d *domain.User) (int64, error) {
	ret := _m.Called(ctx, tx, d)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.User) int64); ok {
		r0 = rf(ctx, tx, d)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *domain.User) error); ok {
		r1 = rf(ctx, tx, d)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateNumberOfNumberOfCollectionTx provides a mock function with given fields: ctx, tx, d
func (_m *UserMySQLRepository) UpdateNumberOfNumberOfCollectionTx(ctx context.Context, tx *gorm.DB, d *domain.User) (int64, error) {
	ret := _m.Called(ctx, tx, d)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.User) int64); ok {
		r0 = rf(ctx, tx, d)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *domain.User) error); ok {
		r1 = rf(ctx, tx, d)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateNumberOfPostTx provides a mock function with given fields: ctx, tx, d
func (_m *UserMySQLRepository) UpdateNumberOfPostTx(ctx context.Context, tx *gorm.DB, d *domain.User) (int64, error) {
	ret := _m.Called(ctx, tx, d)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.User) int64); ok {
		r0 = rf(ctx, tx, d)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *domain.User) error); ok {
		r1 = rf(ctx, tx, d)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
