// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/beecool-cocktail/application-backend/domain"
	mock "github.com/stretchr/testify/mock"
)

// CocktailUsecase is an autogenerated mock type for the CocktailUsecase type
type CocktailUsecase struct {
	mock.Mock
}

// GetAllWithFilter provides a mock function with given fields: ctx, filter, pagination
func (_m *CocktailUsecase) GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination domain.PaginationUsecase) ([]domain.Cocktail, int64, error) {
	ret := _m.Called(ctx, filter, pagination)

	var r0 []domain.Cocktail
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, domain.PaginationUsecase) []domain.Cocktail); ok {
		r0 = rf(ctx, filter, pagination)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Cocktail)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}, domain.PaginationUsecase) int64); ok {
		r1 = rf(ctx, filter, pagination)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, map[string]interface{}, domain.PaginationUsecase) error); ok {
		r2 = rf(ctx, filter, pagination)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
