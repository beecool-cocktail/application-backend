package usercase

import (
	"context"
	"mime/multipart"
	"testing"

	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
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

func Test_userUsecase_UpdateBasicInfo(t *testing.T) {
	mockUserMySQLRepo := new(mocks.UserMySQLRepository)
	mockUserRedisRepo := new(mocks.UserRedisRepository)

	mockUser := domain.User{
		ID:                 1,
		Name:               "Andy",
		IsCollectionPublic: true,
	}

	t.Run("Success", func(t *testing.T) {
		mockUserMySQLRepo.
			On("UpdateBasicInfo", mock.Anything, mock.Anything).
			Return(int64(1), nil).Once()

		mockUserRedisRepo.
			On("UpdateBasicInfo", mock.Anything, mock.Anything).
			Return( nil).Once()

		u := NewUserUsecase(mockUserMySQLRepo, mockUserRedisRepo)
		_, err := u.UpdateBasicInfo(context.TODO(), &mockUser)

		assert.NoError(t, err)
		mockUserMySQLRepo.AssertExpectations(t)
	})
}

func Test_userUsecase_UpdateImage(t *testing.T) {
	mockUserMySQLRepo := new(mocks.UserMySQLRepository)
	mockUserRedisRepo := new(mocks.UserRedisRepository)

	mockUser := domain.UserImage{
		ID:                 1,
		File: &multipart.FileHeader{
			Filename: "filename",
		},
	}

	t.Run("Success", func(t *testing.T) {
		mockUserMySQLRepo.
			On("UpdateImage",
				mock.Anything,
				mock.MatchedBy(func(u *domain.UserImage) bool { return u.Destination != "" })).
			Return(int64(1), nil).Once()

		u := NewUserUsecase(mockUserMySQLRepo, mockUserRedisRepo)
		_, err := u.UpdateImage(context.TODO(), &mockUser)

		assert.NoError(t, err)
		mockUserMySQLRepo.AssertExpectations(t)
	})
}
