// Code generated by mockery v2.30.16. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/beecool-cocktail/application-backend/domain"
	mock "github.com/stretchr/testify/mock"
)

// FavoriteCocktailUsecase is an autogenerated mock type for the FavoriteCocktailUsecase type
type FavoriteCocktailUsecase struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, cocktailID, userID
func (_m *FavoriteCocktailUsecase) Delete(ctx context.Context, cocktailID int64, userID int64) (string, error) {
	ret := _m.Called(ctx, cocktailID, userID)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) (string, error)); ok {
		return rf(ctx, cocktailID, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) string); ok {
		r0 = rf(ctx, cocktailID, userID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, cocktailID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryByUserID provides a mock function with given fields: ctx, targetUserID, pagination, needCollectedStatusUserID
func (_m *FavoriteCocktailUsecase) QueryByUserID(ctx context.Context, targetUserID int64, pagination domain.PaginationUsecase, needCollectedStatusUserID int64) ([]domain.APIFavoriteCocktail, int64, error) {
	ret := _m.Called(ctx, targetUserID, pagination, needCollectedStatusUserID)

	var r0 []domain.APIFavoriteCocktail
	var r1 int64
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, domain.PaginationUsecase, int64) ([]domain.APIFavoriteCocktail, int64, error)); ok {
		return rf(ctx, targetUserID, pagination, needCollectedStatusUserID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, domain.PaginationUsecase, int64) []domain.APIFavoriteCocktail); ok {
		r0 = rf(ctx, targetUserID, pagination, needCollectedStatusUserID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.APIFavoriteCocktail)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, domain.PaginationUsecase, int64) int64); ok {
		r1 = rf(ctx, targetUserID, pagination, needCollectedStatusUserID)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(context.Context, int64, domain.PaginationUsecase, int64) error); ok {
		r2 = rf(ctx, targetUserID, pagination, needCollectedStatusUserID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// QueryCountsByUserID provides a mock function with given fields: ctx, id
func (_m *FavoriteCocktailUsecase) QueryCountsByUserID(ctx context.Context, id int64) (int64, error) {
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

// Store provides a mock function with given fields: ctx, c
func (_m *FavoriteCocktailUsecase) Store(ctx context.Context, c *domain.FavoriteCocktail) error {
	ret := _m.Called(ctx, c)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.FavoriteCocktail) error); ok {
		r0 = rf(ctx, c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewFavoriteCocktailUsecase creates a new instance of FavoriteCocktailUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFavoriteCocktailUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *FavoriteCocktailUsecase {
	mock := &FavoriteCocktailUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
