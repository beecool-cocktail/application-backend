package usercase

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"gorm.io/gorm"
	"testing"

	"github.com/beecool-cocktail/application-backend/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_userUsecase_Logout(t *testing.T) {
	mockUserMySQLRepo := new(mocks.UserMySQLRepository)
	mockUserRedisRepo := new(mocks.UserRedisRepository)

	t.Run("Success", func(t *testing.T) {
		mockUserRedisRepo.
			On("UpdateToken", mock.Anything, mock.Anything).
			Return(nil, nil).Once()

		u := NewUserUsecase(mockUserMySQLRepo, mockUserRedisRepo)
		err := u.Logout(context.TODO(), 1)

		assert.NoError(t, err)
	})
}

func Test_userUsecase_QueryById(t *testing.T) {
	mockUserMySQLRepo := new(mocks.UserMySQLRepository)
	mockUserRedisRepo := new(mocks.UserRedisRepository)

	mockUser := domain.User{
		ID: 1,
	}

	t.Run("Success", func(t *testing.T) {
		id := int64(1)
		mockUserMySQLRepo.
			On("QueryById", mock.Anything, id).
			Return(&mockUser, nil).Once()

		u := NewUserUsecase(mockUserMySQLRepo, mockUserRedisRepo)
		user, err := u.QueryById(context.TODO(), 1)

		assert.NoError(t, err)
		assert.Equal(t, &mockUser, user)
		mockUserMySQLRepo.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		id := int64(1)
		mockUserMySQLRepo.
			On("QueryById", mock.Anything, id).
			Return(&mockUser, gorm.ErrRecordNotFound).Once()

		u := NewUserUsecase(mockUserMySQLRepo, mockUserRedisRepo)
		_, err := u.QueryById(context.TODO(), 1)

		assert.Equal(t, domain.ErrUserNotFound, err)
		mockUserMySQLRepo.AssertExpectations(t)
	})
}
