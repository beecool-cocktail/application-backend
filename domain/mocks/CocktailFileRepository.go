// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/beecool-cocktail/application-backend/domain"
	mock "github.com/stretchr/testify/mock"
)

// CocktailFileRepository is an autogenerated mock type for the CocktailFileRepository type
type CocktailFileRepository struct {
	mock.Mock
}

// SaveAsWebp provides a mock function with given fields: ctx, ci
func (_m *CocktailFileRepository) SaveAsWebp(ctx context.Context, ci *domain.CocktailImage) error {
	ret := _m.Called(ctx, ci)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.CocktailImage) error); ok {
		r0 = rf(ctx, ci)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveAsWebpInLQIP provides a mock function with given fields: ctx, ci
func (_m *CocktailFileRepository) SaveAsWebpInLQIP(ctx context.Context, ci *domain.CocktailImage) error {
	ret := _m.Called(ctx, ci)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.CocktailImage) error); ok {
		r0 = rf(ctx, ci)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateAsWebp provides a mock function with given fields: ctx, ci
func (_m *CocktailFileRepository) UpdateAsWebp(ctx context.Context, ci *domain.CocktailImage) error {
	ret := _m.Called(ctx, ci)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.CocktailImage) error); ok {
		r0 = rf(ctx, ci)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateAsWebpInLQIP provides a mock function with given fields: ctx, ci
func (_m *CocktailFileRepository) UpdateAsWebpInLQIP(ctx context.Context, ci *domain.CocktailImage) error {
	ret := _m.Called(ctx, ci)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.CocktailImage) error); ok {
		r0 = rf(ctx, ci)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
