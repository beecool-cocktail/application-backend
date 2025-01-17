// Code generated by mockery v2.30.16. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/beecool-cocktail/application-backend/domain"
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

// FavoriteCocktailMySQLRepository is an autogenerated mock type for the FavoriteCocktailMySQLRepository type
type FavoriteCocktailMySQLRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, cocktailID, userID
func (_m *FavoriteCocktailMySQLRepository) Delete(ctx context.Context, cocktailID int64, userID int64) error {
	ret := _m.Called(ctx, cocktailID, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, cocktailID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTx provides a mock function with given fields: ctx, tx, cocktailID, userID
func (_m *FavoriteCocktailMySQLRepository) DeleteTx(ctx context.Context, tx *gorm.DB, cocktailID int64, userID int64) error {
	ret := _m.Called(ctx, tx, cocktailID, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, int64, int64) error); ok {
		r0 = rf(ctx, tx, cocktailID, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueryByUserID provides a mock function with given fields: ctx, id, pagination
func (_m *FavoriteCocktailMySQLRepository) QueryByUserID(ctx context.Context, id int64, pagination domain.PaginationMySQLRepository) ([]domain.FavoriteCocktail, int64, error) {
	ret := _m.Called(ctx, id, pagination)

	var r0 []domain.FavoriteCocktail
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, domain.PaginationMySQLRepository) ([]domain.FavoriteCocktail, int64, error)); ok {
		return rf(ctx, id, pagination)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, domain.PaginationMySQLRepository) []domain.FavoriteCocktail); ok {
		r0 = rf(ctx, id, pagination)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.FavoriteCocktail)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, domain.PaginationMySQLRepository) int64); ok {
		r1 = rf(ctx, id, pagination)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, int64, domain.PaginationMySQLRepository) error); ok {
		r2 = rf(ctx, id, pagination)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// QueryCountsByUserID provides a mock function with given fields: ctx, id
func (_m *FavoriteCocktailMySQLRepository) QueryCountsByUserID(ctx context.Context, id int64) (int64, error) {
	ret := _m.Called(ctx, id)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (int64, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) int64); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreTx provides a mock function with given fields: ctx, tx, c
func (_m *FavoriteCocktailMySQLRepository) StoreTx(ctx context.Context, tx *gorm.DB, c *domain.FavoriteCocktail) error {
	ret := _m.Called(ctx, tx, c)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.FavoriteCocktail) error); ok {
		r0 = rf(ctx, tx, c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewFavoriteCocktailMySQLRepository creates a new instance of FavoriteCocktailMySQLRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFavoriteCocktailMySQLRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *FavoriteCocktailMySQLRepository {
	mock := &FavoriteCocktailMySQLRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
