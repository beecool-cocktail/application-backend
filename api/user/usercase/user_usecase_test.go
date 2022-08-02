package usercase

import (
	"context"
	"errors"
	_dbRepo "github.com/beecool-cocktail/application-backend/db/repository/mysql"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/domain/mocks"
	"github.com/beecool-cocktail/application-backend/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_userUsecase_Logout(t *testing.T) {
	db, _, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	svc := testutil.GetService()

	mockUserMySQLRepo := new(mocks.UserMySQLRepository)
	mockUserRedisRepo := new(mocks.UserRedisRepository)
	mockUserFileRepo := new(mocks.UserFileRepository)
	mockTransactionMySQLRepo := _dbRepo.NewDBRepository(db)

	t.Run("Success", func(t *testing.T) {
		id := int64(123456)

		mockUserRedisRepo.
			On("UpdateToken",
				mock.Anything,
				mock.MatchedBy(func(token *domain.UserCache) bool {
					return matchedByUpdateToken(token, id)
				})).
			Return(nil).Once()

		u := NewUserUsecase(svc, mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionMySQLRepo)
		err := u.Logout(context.TODO(), id)

		assert.NoError(t, err)

		mockUserRedisRepo.AssertExpectations(t)
	})
}

func Test_userUsecase_QueryById(t *testing.T) {
	db, _, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	svc := testutil.GetService()

	mockUserMySQLRepo := new(mocks.UserMySQLRepository)
	mockUserRedisRepo := new(mocks.UserRedisRepository)
	mockUserFileRepo := new(mocks.UserFileRepository)
	mockTransactionMySQLRepo := _dbRepo.NewDBRepository(db)

	t.Run("Success", func(t *testing.T) {
		mockUser := domain.User{
			ID:                 123456,
			Account:            "account",
			Password:           "password",
			Status:             0,
			Name:               "name",
			Email:              "email",
			Photo:              "photo",
			NumberOfPost:       1,
			NumberOfCollection: 2,
			NumberOfDraft:      3,
			IsCollectionPublic: true,
			Remark:             "good",
			CreatedDate:        time.Time{},
		}
		id := int64(123456)

		mockUserMySQLRepo.
			On("QueryById",
				mock.Anything,
				id).
			Return(mockUser, nil).Once()

		u := NewUserUsecase(svc, mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionMySQLRepo)
		user, err := u.QueryById(context.TODO(), id)

		assert.NoError(t, err)
		assert.Equal(t, mockUser, user)

		mockUserRedisRepo.AssertExpectations(t)
	})
}

func Test_userUsecase_UpdateUserInfo(t *testing.T) {
	db, mockDB, err := testutil.BeforeEach()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	svc := testutil.GetService()

	t.Run("Success - doesn't update image", func(t *testing.T) {
		mockUserMySQLRepo := new(mocks.UserMySQLRepository)
		mockUserRedisRepo := new(mocks.UserRedisRepository)
		mockUserFileRepo := new(mocks.UserFileRepository)
		mockTransactionMySQLRepo := _dbRepo.NewDBRepository(db)
		mockUser := domain.User{
			ID:                 123456,
			Account:            "account",
			Password:           "password",
			Status:             0,
			Name:               "name",
			Email:              "email",
			Photo:              "photo",
			Height:             0,
			Width:              0,
			NumberOfPost:       1,
			NumberOfCollection: 2,
			NumberOfDraft:      3,
			IsCollectionPublic: true,
			Remark:             "good",
			CreatedDate:        time.Time{},
		}

		mockUserCache := domain.UserCache{
			Name: mockUser.Name,
		}

		mockUserImage := domain.UserImage{
			ID:   123456,
			Data: "",
		}

		mockDB.ExpectBegin()

		mockUserMySQLRepo.
			On("UpdateBasicInfoTx",
				mock.Anything,
				mock.Anything,
				&mockUser).
			Return(int64(1), nil).Once()

		mockUserRedisRepo.
			On("UpdateBasicInfo",
				mock.Anything,
				&mockUserCache).
			Return(nil).Once()

		mockDB.ExpectCommit()

		u := NewUserUsecase(svc, mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionMySQLRepo)
		err := u.UpdateUserInfo(context.TODO(), &mockUser, &mockUserImage)

		assert.NoError(t, err)

		mockUserMySQLRepo.AssertExpectations(t)
		mockUserRedisRepo.AssertExpectations(t)
	})

	t.Run("Success - need update image", func(t *testing.T) {
		mockUserMySQLRepo := new(mocks.UserMySQLRepository)
		mockUserRedisRepo := new(mocks.UserRedisRepository)
		mockUserFileRepo := new(mocks.UserFileRepository)
		mockTransactionMySQLRepo := _dbRepo.NewDBRepository(db)
		mockUser := domain.User{
			ID:                 123456,
			Account:            "account",
			Password:           "password",
			Status:             0,
			Name:               "name",
			Email:              "email",
			Photo:              "photo",
			Width:              10,
			Height:             20,
			NumberOfPost:       1,
			NumberOfCollection: 2,
			NumberOfDraft:      3,
			IsCollectionPublic: true,
			Remark:             "good",
			CreatedDate:        time.Time{},
		}

		mockUserCache := domain.UserCache{
			Name: mockUser.Name,
		}

		mockFileName := "file name"

		mockUserImage := domain.UserImage{
			ID:   123456,
			Data: "dataURL",
			Name: mockFileName,
		}

		mockDB.ExpectBegin()

		mockUserImage.Destination = svc.Configure.Others.File.Image.PathInDB + mockFileName
		mockUserFileRepo.
			On("SaveAsWebp",
				mock.Anything,
				&mockUserImage).
			Return(10, 20, nil).Once()

		mockUserImage.Destination = svc.Configure.Others.File.Image.PathInURL + mockFileName + ".webp"
		mockUserMySQLRepo.
			On("UpdateImageTx",
				mock.Anything,
				mock.Anything,
				&mockUserImage).
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
				&mockUserCache).
			Return(nil).Once()

		mockDB.ExpectCommit()

		u := NewUserUsecase(svc, mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionMySQLRepo)
		err := u.UpdateUserInfo(context.TODO(), &mockUser, &mockUserImage)

		assert.NoError(t, err)

		mockUserMySQLRepo.AssertExpectations(t)
		mockUserRedisRepo.AssertExpectations(t)
		mockUserFileRepo.AssertExpectations(t)
	})

	t.Run("Failed - roll back", func(t *testing.T) {
		mockUserMySQLRepo := new(mocks.UserMySQLRepository)
		mockUserRedisRepo := new(mocks.UserRedisRepository)
		mockUserFileRepo := new(mocks.UserFileRepository)
		mockTransactionMySQLRepo := _dbRepo.NewDBRepository(db)
		mockUser := domain.User{
			ID:                 123456,
			Account:            "account",
			Password:           "password",
			Status:             0,
			Name:               "name",
			Email:              "email",
			Photo:              "photo",
			NumberOfPost:       1,
			NumberOfCollection: 2,
			NumberOfDraft:      3,
			IsCollectionPublic: true,
			Remark:             "good",
			CreatedDate:        time.Time{},
		}

		mockUserCache := domain.UserCache{
			Name: mockUser.Name,
		}

		mockUserImage := domain.UserImage{
			ID:   123456,
			Data: "",
		}

		mockDB.ExpectBegin()

		mockUserMySQLRepo.
			On("UpdateBasicInfoTx",
				mock.Anything,
				mock.Anything,
				&mockUser).
			Return(int64(1), nil).Once()

		mockUserRedisRepo.
			On("UpdateBasicInfo",
				mock.Anything,
				&mockUserCache).
			Return(errors.New("failed")).Once()

		mockDB.ExpectRollback()

		u := NewUserUsecase(svc, mockUserMySQLRepo, mockUserRedisRepo, mockUserFileRepo, mockTransactionMySQLRepo)
		err := u.UpdateUserInfo(context.TODO(), &mockUser, &mockUserImage)

		assert.EqualError(t, err, errors.New("failed").Error())

		mockUserMySQLRepo.AssertExpectations(t)
		mockUserRedisRepo.AssertExpectations(t)
	})
}
