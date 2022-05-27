package usecase

import (
	"context"
	"errors"
	"github.com/beecool-cocktail/application-backend/domain"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"testing"

	"github.com/beecool-cocktail/application-backend/domain/mocks"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_socialLoginUsecase_Exchange(t *testing.T) {
	mockUserMySQLRepo := new(mocks.UserMySQLRepository)
	mockUserRedisRepo := new(mocks.UserRedisRepository)
	mockSocialAccountGoogleOAuth2Repo := new(mocks.SocialAccountGoogleOAuthRepository)
	mockSocialAccountMySQLRepo := new(mocks.SocialAccountMySQLRepository)

	go util.StartUserIdGenerator()

	t.Run("Success", func(t *testing.T) {
		mockSocialAccountGoogleOAuth2Repo.
			On("Exchange", mock.Anything, mock.MatchedBy(func(code string) bool { return code == "code" })).
			Return(nil, nil).Once()

		s := NewSocialAccountUsecase(mockUserMySQLRepo, mockUserRedisRepo, mockSocialAccountMySQLRepo, mockSocialAccountGoogleOAuth2Repo)
		_, err := s.Exchange(context.TODO(), "code")

		assert.NoError(t, err)
	})
}

func Test_socialLoginUsecase_GetUserInfo(t *testing.T) {
	mockUserMySQLRepo := new(mocks.UserMySQLRepository)
	mockUserRedisRepo := new(mocks.UserRedisRepository)
	mockSocialAccountGoogleOAuth2Repo := new(mocks.SocialAccountGoogleOAuthRepository)
	mockSocialAccountMySQLRepo := new(mocks.SocialAccountMySQLRepository)

	oauthToken := &oauth2.Token{
		AccessToken: "token",
	}

	mockSocialAccount := domain.SocialAccount{
		UserID:   1,
		SocialID: "googleUUID",
	}

	mockGoogleUserInfo := domain.GoogleUserInfo{
		Name: "Andy",
		Sub:  "googleUUID",
	}

	mockUserMySQL := domain.User{
		ID:      1,
		Account: "account",
		Name:    "Andy",
	}

	t.Run("Not first login Success", func(t *testing.T) {
		mockSocialAccountGoogleOAuth2Repo.
			On("GetUserInfo", mock.Anything, mock.MatchedBy(func(token *oauth2.Token) bool { return token == oauthToken })).
			Return(&mockGoogleUserInfo, nil).Once()

		mockSocialAccountMySQLRepo.
			On("QueryById", mock.Anything, mock.MatchedBy(func(googleUUID string) bool { return googleUUID == mockGoogleUserInfo.Sub })).
			Return(&mockSocialAccount, nil).Once()

		mockUserMySQLRepo.
			On("QueryById", mock.Anything, mock.MatchedBy(func(userID int64) bool { return userID == mockSocialAccount.UserID })).
			Return(&mockUserMySQL, nil).Once()

		mockUserRedisRepo.
			On("Store", mock.Anything, mock.MatchedBy(func(mockUser *domain.UserCache) bool {
				return matchedByUserRedis(mockUser, &mockUserMySQL)
			})).
			Return(nil).Once()

		s := NewSocialAccountUsecase(mockUserMySQLRepo, mockUserRedisRepo, mockSocialAccountMySQLRepo, mockSocialAccountGoogleOAuth2Repo)
		_, err := s.GetUserInfo(context.TODO(), oauthToken)

		assert.NoError(t, err)
	})

	t.Run("First login Success", func(t *testing.T) {
		mockSocialAccountGoogleOAuth2Repo.
			On("GetUserInfo", mock.Anything, mock.MatchedBy(func(token *oauth2.Token) bool { return token == oauthToken })).
			Return(&mockGoogleUserInfo, nil).Once()

		mockSocialAccountMySQLRepo.
			On("QueryById", mock.Anything, mock.MatchedBy(func(googleUUID string) bool { return googleUUID == mockGoogleUserInfo.Sub })).
			Return(nil, gorm.ErrRecordNotFound).Once()

		mockSocialAccountMySQLRepo.
			On("Store", mock.Anything, mock.Anything, mock.Anything).
			Return(int64(1), nil).Once()

		mockUserRedisRepo.
			On("Store", mock.Anything, mock.Anything).
			Return(nil).Once()

		s := NewSocialAccountUsecase(mockUserMySQLRepo, mockUserRedisRepo, mockSocialAccountMySQLRepo, mockSocialAccountGoogleOAuth2Repo)
		_, err := s.GetUserInfo(context.TODO(), oauthToken)

		assert.NoError(t, err)
	})

	t.Run("Get google user info Failed", func(t *testing.T) {
		mockSocialAccountGoogleOAuth2Repo.
			On("GetUserInfo", mock.Anything, mock.MatchedBy(func(token *oauth2.Token) bool { return token == oauthToken })).
			Return(nil, errors.New("get google user info failed")).Once()

		s := NewSocialAccountUsecase(mockUserMySQLRepo, mockUserRedisRepo, mockSocialAccountMySQLRepo, mockSocialAccountGoogleOAuth2Repo)
		_, err := s.GetUserInfo(context.TODO(), oauthToken)

		assert.Equal(t, err, errors.New("get google user info failed"))
	})

	t.Run("Get social account Failed", func(t *testing.T) {
		mockSocialAccountGoogleOAuth2Repo.
			On("GetUserInfo", mock.Anything, mock.MatchedBy(func(token *oauth2.Token) bool { return token == oauthToken })).
			Return(&mockGoogleUserInfo, nil).Once()

		mockSocialAccountMySQLRepo.
			On("QueryById", mock.Anything, mock.MatchedBy(func(googleUUID string) bool { return googleUUID == mockGoogleUserInfo.Sub })).
			Return(nil, errors.New("get social account info failed")).Once()

		s := NewSocialAccountUsecase(mockUserMySQLRepo, mockUserRedisRepo, mockSocialAccountMySQLRepo, mockSocialAccountGoogleOAuth2Repo)
		_, err := s.GetUserInfo(context.TODO(), oauthToken)

		assert.Equal(t, err, errors.New("get social account info failed"))
	})
}
