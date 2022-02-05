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
	mockUserFileRepo := new(mocks.UserFileRepository)
	mockTransactionRepo := new(mocks.DBTransactionRepository)

	t.Run("Success", func(t *testing.T) {
		mockUserRedisRepo.
			On("UpdateToken", mock.Anything, mock.Anything).
			Return(nil, nil).Once()

		u := NewUserUsecase(mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionRepo)
		err := u.Logout(context.TODO(), 1)

		assert.NoError(t, err)
	})
}

func Test_userUsecase_QueryById(t *testing.T) {
	mockUserMySQLRepo := new(mocks.UserMySQLRepository)
	mockUserRedisRepo := new(mocks.UserRedisRepository)
	mockUserFileRepo := new(mocks.UserFileRepository)
	mockTransactionRepo := new(mocks.DBTransactionRepository)

	mockUser := domain.User{
		ID: 1,
	}

	t.Run("Success", func(t *testing.T) {
		id := int64(1)
		mockUserMySQLRepo.
			On("QueryById", mock.Anything, id).
			Return(&mockUser, nil).Once()

		u := NewUserUsecase(mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionRepo)
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

		u := NewUserUsecase(mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionRepo)
		_, err := u.QueryById(context.TODO(), 1)

		assert.Equal(t, domain.ErrUserNotFound, err)
		mockUserMySQLRepo.AssertExpectations(t)
	})
}

func Test_userUsecase_UpdateBasicInfo(t *testing.T) {
	mockUserMySQLRepo := new(mocks.UserMySQLRepository)
	mockUserRedisRepo := new(mocks.UserRedisRepository)
	mockUserFileRepo := new(mocks.UserFileRepository)
	mockTransactionRepo := new(mocks.DBTransactionRepository)

	mockUser := domain.User{
		ID:                 1,
		Name:               "Andy",
		IsCollectionPublic: true,
	}

	mockUserImage := domain.UserImage{
		ID:                 1,
		Data: &multipart.FileHeader{
			Filename: "filename.png",
		},
		Type: ".png",
	}

	t.Run("Success", func(t *testing.T) {
		mockTransactionRepo.
			On("Transaction", mock.Anything).
			Return(nil).Once()
		
		mockUserRedisRepo.
			On("UpdateBasicInfo", mock.Anything, mock.Anything).
			Return( nil).Once()

		u := NewUserUsecase(mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionRepo)
		err := u.UpdateUserInfo(context.TODO(), &mockUser, &mockUserImage)

		assert.NoError(t, err)
		mockUserMySQLRepo.AssertExpectations(t)
	})
}
