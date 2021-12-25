package usercase

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
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