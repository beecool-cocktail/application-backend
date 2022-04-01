package usercase

import (
	"context"
	"errors"
	"github.com/beecool-cocktail/application-backend/testutil"
	"testing"

	_transactionRepo "github.com/beecool-cocktail/application-backend/db/repository/mysql"
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

	mockService := testutil.GetService()

	t.Run("Success", func(t *testing.T) {
		mockUserRedisRepo.
			On("UpdateToken", mock.Anything, mock.Anything).
			Return(nil, nil).Once()

		u := NewUserUsecase(mockService, mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionRepo)
		err := u.Logout(context.TODO(), 1)

		assert.NoError(t, err)
	})
}

func Test_userUsecase_QueryById(t *testing.T) {
	mockUserMySQLRepo := new(mocks.UserMySQLRepository)
	mockUserRedisRepo := new(mocks.UserRedisRepository)
	mockUserFileRepo := new(mocks.UserFileRepository)
	mockTransactionRepo := new(mocks.DBTransactionRepository)

	mockService := testutil.GetService()

	mockUser := domain.User{
		ID: 1,
	}

	t.Run("Success", func(t *testing.T) {
		id := int64(1)
		mockUserMySQLRepo.
			On("QueryById", mock.Anything, id).
			Return(&mockUser, nil).Once()

		u := NewUserUsecase(mockService, mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionRepo)
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

		u := NewUserUsecase(mockService, mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionRepo)
		_, err := u.QueryById(context.TODO(), 1)

		assert.Equal(t, domain.ErrUserNotFound, err)
		mockUserMySQLRepo.AssertExpectations(t)
	})
}

func Test_userUsecase_UpdateUserInfo(t *testing.T) {
	db, dbMock, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	mockUserMySQLRepo := new(mocks.UserMySQLRepository)
	mockUserRedisRepo := new(mocks.UserRedisRepository)
	mockUserFileRepo := new(mocks.UserFileRepository)
	mockTransactionRepo := _transactionRepo.NewDBRepository(db)

	mockService := testutil.GetService()

	t.Run("Success - didn't update image", func(t *testing.T) {
		mockUser := domain.User{
			ID:                 1,
			Name:               "Andy",
			IsCollectionPublic: true,
		}

		mockUserImage := domain.UserImage{
			ID:   1,
			Type: "image/png",
		}

		dbMock.ExpectBegin()
		mockUserMySQLRepo.
			On("UpdateBasicInfoTx",
				mock.Anything,
				mock.Anything,
				&mockUser).
			Return(int64(1), nil).Once()

		mockUserRedisRepo.
			On("UpdateBasicInfo",
				mock.Anything,
				mock.MatchedBy(func(userCache *domain.UserCache) bool {
					return matchedByUserOfUpdateUserInfo(userCache, &mockUser)
				})).
			Return(nil).Once()

		dbMock.ExpectCommit()

		u := NewUserUsecase(mockService, mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionRepo)
		err := u.UpdateUserInfo(context.TODO(), &mockUser, &mockUserImage)

		assert.NoError(t, err)
		mockUserMySQLRepo.AssertExpectations(t)
	})

	t.Run("Success - update image", func(t *testing.T) {
		mockUser := domain.User{
			ID:                 1,
			Name:               "Andy",
			IsCollectionPublic: true,
		}

		mockUserImage := domain.UserImage{
			ID:   1,
			Data: "path",
			Type: "image/png",
		}

		dbMock.ExpectBegin()
		mockUserFileRepo.
			On("SaveAsWebp",
				mock.Anything,
				mock.Anything).
			Return(nil).Once()

		mockUserMySQLRepo.
			On("UpdateImageTx",
				mock.Anything,
				mock.Anything,
				mock.Anything).
			Return(int64(1), nil).Once()

		mockUserMySQLRepo.
			On("UpdateBasicInfoTx",
				mock.Anything,
				mock.Anything,
				&mockUser).
			Return(int64(1), nil).Once()

		mockUserRedisRepo.
			On("UpdateBasicInfo",
				mock.Anything,
				mock.MatchedBy(func(userCache *domain.UserCache) bool {
					return matchedByUserOfUpdateUserInfo(userCache, &mockUser)
				})).
			Return(nil).Once()

		dbMock.ExpectCommit()

		u := NewUserUsecase(mockService, mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionRepo)
		err := u.UpdateUserInfo(context.TODO(), &mockUser, &mockUserImage)

		assert.NoError(t, err)
		mockUserMySQLRepo.AssertExpectations(t)
	})

	t.Run("Success - didn't update image", func(t *testing.T) {
		mockUser := domain.User{
			ID:                 1,
			Name:               "Andy",
			IsCollectionPublic: true,
		}

		mockUserImage := domain.UserImage{
			ID:   1,
			Type: "image/png",
		}

		dbMock.ExpectBegin()
		mockUserMySQLRepo.
			On("UpdateBasicInfoTx",
				mock.Anything,
				mock.Anything,
				&mockUser).
			Return(int64(1), nil).Once()

		mockUserRedisRepo.
			On("UpdateBasicInfo",
				mock.Anything,
				mock.MatchedBy(func(userCache *domain.UserCache) bool {
					return matchedByUserOfUpdateUserInfo(userCache, &mockUser)
				})).
			Return(errors.New("update failed")).Once()

		dbMock.ExpectRollback()

		u := NewUserUsecase(mockService, mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionRepo)
		err := u.UpdateUserInfo(context.TODO(), &mockUser, &mockUserImage)

		assert.EqualError(t, err, "update failed")
		mockUserMySQLRepo.AssertExpectations(t)
	})
}
