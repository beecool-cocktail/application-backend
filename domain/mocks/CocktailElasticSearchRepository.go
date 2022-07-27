// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/beecool-cocktail/application-backend/domain"
	mock "github.com/stretchr/testify/mock"
)

// CocktailElasticSearchRepository is an autogenerated mock type for the CocktailElasticSearchRepository type
type CocktailElasticSearchRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, id
func (_m *CocktailElasticSearchRepository) Delete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Index provides a mock function with given fields: ctx, c
func (_m *CocktailElasticSearchRepository) Index(ctx context.Context, c *domain.CocktailElasticSearch) error {
	ret := _m.Called(ctx, c)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.CocktailElasticSearch) error); ok {
		r0 = rf(ctx, c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Search provides a mock function with given fields: ctx, text, from, size
func (_m *CocktailElasticSearchRepository) Search(ctx context.Context, text string, from int, size int) ([]domain.CocktailElasticSearch, int64, error) {
	ret := _m.Called(ctx, text, from, size)

	var r0 []domain.CocktailElasticSearch
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []domain.CocktailElasticSearch); ok {
		r0 = rf(ctx, text, from, size)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.CocktailElasticSearch)
		}
	}

	var r1 int64
	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) int64); ok {
		r1 = rf(ctx, text, from, size)
	} else {
		r1 = ret.Get(1).(int64)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, int, int) error); ok {
		r2 = rf(ctx, text, from, size)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Update provides a mock function with given fields: ctx, c
func (_m *CocktailElasticSearchRepository) Update(ctx context.Context, c *domain.CocktailElasticSearch) error {
	ret := _m.Called(ctx, c)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.CocktailElasticSearch) error); ok {
		r0 = rf(ctx, c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}