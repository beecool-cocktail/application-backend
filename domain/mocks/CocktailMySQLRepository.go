// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/beecool-cocktail/application-backend/domain"
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

// CocktailMySQLRepository is an autogenerated mock type for the CocktailMySQLRepository type
type CocktailMySQLRepository struct {
	mock.Mock
}

// DeleteTx provides a mock function with given fields: ctx, tx, id
func (_m *CocktailMySQLRepository) DeleteTx(ctx context.Context, tx *gorm.DB, id int64) error {
	ret := _m.Called(ctx, tx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, int64) error); ok {
		r0 = rf(ctx, tx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllWithFilter provides a mock function with given fields: ctx, filter, pagination
func (_m *CocktailMySQLRepository) GetAllWithFilter(ctx context.Context, filter map[string]interface{}, pagination domain.PaginationMySQLRepository) ([]domain.Cocktail, int64, error) {
	ret := _m.Called(ctx, filter, pagination)

	var r0 []domain.Cocktail
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, domain.PaginationMySQLRepository) []domain.Cocktail); ok {
		r0 = rf(ctx, filter, pagination)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Cocktail)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}, domain.PaginationMySQLRepository) int64); ok {
		r1 = rf(ctx, filter, pagination)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, map[string]interface{}, domain.PaginationMySQLRepository) error); ok {
		r2 = rf(ctx, filter, pagination)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// QueryByCocktailID provides a mock function with given fields: ctx, id
func (_m *CocktailMySQLRepository) QueryByCocktailID(ctx context.Context, id int64) (domain.Cocktail, error) {
	ret := _m.Called(ctx, id)

	var r0 domain.Cocktail
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.Cocktail); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.Cocktail)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreTx provides a mock function with given fields: ctx, tx, c
func (_m *CocktailMySQLRepository) StoreTx(ctx context.Context, tx *gorm.DB, c *domain.Cocktail) error {
	ret := _m.Called(ctx, tx, c)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.Cocktail) error); ok {
		r0 = rf(ctx, tx, c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateCategoryTx provides a mock function with given fields: ctx, tx, c
func (_m *CocktailMySQLRepository) UpdateCategoryTx(ctx context.Context, tx *gorm.DB, c *domain.Cocktail) (int64, error) {
	ret := _m.Called(ctx, tx, c)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.Cocktail) int64); ok {
		r0 = rf(ctx, tx, c)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *domain.Cocktail) error); ok {
		r1 = rf(ctx, tx, c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTx provides a mock function with given fields: ctx, tx, c
func (_m *CocktailMySQLRepository) UpdateTx(ctx context.Context, tx *gorm.DB, c *domain.Cocktail) (int64, error) {
	ret := _m.Called(ctx, tx, c)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.Cocktail) int64); ok {
		r0 = rf(ctx, tx, c)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *domain.Cocktail) error); ok {
		r1 = rf(ctx, tx, c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
