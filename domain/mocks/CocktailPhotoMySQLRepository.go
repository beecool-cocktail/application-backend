// Code generated by mockery v2.30.16. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/beecool-cocktail/application-backend/domain"
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

// CocktailPhotoMySQLRepository is an autogenerated mock type for the CocktailPhotoMySQLRepository type
type CocktailPhotoMySQLRepository struct {
	mock.Mock
}

// DeleteByCocktailIDTx provides a mock function with given fields: ctx, tx, id
func (_m *CocktailPhotoMySQLRepository) DeleteByCocktailIDTx(ctx context.Context, tx *gorm.DB, id int64) error {
	ret := _m.Called(ctx, tx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, int64) error); ok {
		r0 = rf(ctx, tx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteByIDTx provides a mock function with given fields: ctx, tx, id
func (_m *CocktailPhotoMySQLRepository) DeleteByIDTx(ctx context.Context, tx *gorm.DB, id int64) error {
	ret := _m.Called(ctx, tx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, int64) error); ok {
		r0 = rf(ctx, tx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteByLowQualityBundleIDTx provides a mock function with given fields: ctx, tx, id
func (_m *CocktailPhotoMySQLRepository) DeleteByLowQualityBundleIDTx(ctx context.Context, tx *gorm.DB, id int64) error {
	ret := _m.Called(ctx, tx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, int64) error); ok {
		r0 = rf(ctx, tx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueryCoverPhotoByCocktailId provides a mock function with given fields: ctx, id
func (_m *CocktailPhotoMySQLRepository) QueryCoverPhotoByCocktailId(ctx context.Context, id int64) (string, error) {
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

// QueryLowQualityPhotoByBundleId provides a mock function with given fields: ctx, id
func (_m *CocktailPhotoMySQLRepository) QueryLowQualityPhotoByBundleId(ctx context.Context, id int64) (domain.CocktailPhoto, error) {
	ret := _m.Called(ctx, id)

	var r0 domain.CocktailPhoto
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (domain.CocktailPhoto, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.CocktailPhoto); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.CocktailPhoto)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryLowQualityPhotosByCocktailId provides a mock function with given fields: ctx, id
func (_m *CocktailPhotoMySQLRepository) QueryLowQualityPhotosByCocktailId(ctx context.Context, id int64) ([]domain.CocktailPhoto, error) {
	ret := _m.Called(ctx, id)

	var r0 []domain.CocktailPhoto
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]domain.CocktailPhoto, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []domain.CocktailPhoto); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.CocktailPhoto)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryPhotoById provides a mock function with given fields: ctx, id
func (_m *CocktailPhotoMySQLRepository) QueryPhotoById(ctx context.Context, id int64) (domain.CocktailPhoto, error) {
	ret := _m.Called(ctx, id)

	var r0 domain.CocktailPhoto
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (domain.CocktailPhoto, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.CocktailPhoto); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.CocktailPhoto)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryPhotosByCocktailId provides a mock function with given fields: ctx, id
func (_m *CocktailPhotoMySQLRepository) QueryPhotosByCocktailId(ctx context.Context, id int64) ([]domain.CocktailPhoto, error) {
	ret := _m.Called(ctx, id)

	var r0 []domain.CocktailPhoto
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]domain.CocktailPhoto, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []domain.CocktailPhoto); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.CocktailPhoto)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreTx provides a mock function with given fields: ctx, tx, c
func (_m *CocktailPhotoMySQLRepository) StoreTx(ctx context.Context, tx *gorm.DB, c *domain.CocktailPhoto) error {
	ret := _m.Called(ctx, tx, c)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.CocktailPhoto) error); ok {
		r0 = rf(ctx, tx, c)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePhotoOrderTx provides a mock function with given fields: ctx, tx, c
func (_m *CocktailPhotoMySQLRepository) UpdatePhotoOrderTx(ctx context.Context, tx *gorm.DB, c *domain.CocktailPhoto) (int64, error) {
	ret := _m.Called(ctx, tx, c)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.CocktailPhoto) (int64, error)); ok {
		return rf(ctx, tx, c)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.CocktailPhoto) int64); ok {
		r0 = rf(ctx, tx, c)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *domain.CocktailPhoto) error); ok {
		r1 = rf(ctx, tx, c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTx provides a mock function with given fields: ctx, tx, c
func (_m *CocktailPhotoMySQLRepository) UpdateTx(ctx context.Context, tx *gorm.DB, c *domain.CocktailPhoto) (int64, error) {
	ret := _m.Called(ctx, tx, c)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.CocktailPhoto) (int64, error)); ok {
		return rf(ctx, tx, c)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *domain.CocktailPhoto) int64); ok {
		r0 = rf(ctx, tx, c)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, *domain.CocktailPhoto) error); ok {
		r1 = rf(ctx, tx, c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCocktailPhotoMySQLRepository creates a new instance of CocktailPhotoMySQLRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCocktailPhotoMySQLRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *CocktailPhotoMySQLRepository {
	mock := &CocktailPhotoMySQLRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
