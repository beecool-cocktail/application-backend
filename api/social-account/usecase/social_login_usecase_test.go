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
	mockUserNySQLRepo := new(mocks.UserMySQLRepository)
	mockSocialAccountGoogleOAuth2Repo := new(mocks.SocialAccountGoogleOAuthRepository)
	mockSocialAccountMySQLRepo := new(mocks.SocialAccountMySQLRepository)

	go util.StartUserIdGenerator()

	t.Run("Success", func(t *testing.T) {
		mockSocialAccountGoogleOAuth2Repo.
			On("Exchange", mock.Anything, mock.MatchedBy(func(code string) bool { return code == "code" })).
			Return(nil, nil).Once()

		s := NewSocialAccountUsecase(mockUserNySQLRepo, mockSocialAccountMySQLRepo, mockSocialAccountGoogleOAuth2Repo)
		_, err := s.Exchange(context.TODO(), "code")

		assert.NoError(t, err)
	})
}

func Test_socialLoginUsecase_GetUserInfo(t *testing.T) {
	mockUserNySQLRepo := new(mocks.UserMySQLRepository)
	mockSocialAccountGoogleOAuth2Repo := new(mocks.SocialAccountGoogleOAuthRepository)
	mockSocialAccountMySQLRepo := new(mocks.SocialAccountMySQLRepository)

	go util.StartUserIdGenerator()

	oauthToken := &oauth2.Token{
		AccessToken: "token",
	}

	mockSocialAccount := domain.SocialAccount{
		SocialID: "googleUUID",
		UserID: 123456,
	}

	mockGoogleUserInfo := domain.GoogleUserInfo{
		ID: "googleUUID",
	}

	mockUser := domain.User{
		UserID: 123456,
		Account: "account",
		Name: "Andy",
	}

	t.Run("First login Success", func(t *testing.T) {
		mockSocialAccountGoogleOAuth2Repo.
			On("GetUserInfo", mock.Anything, mock.MatchedBy(func(token *oauth2.Token) bool { return token == oauthToken })).
			Return(&mockGoogleUserInfo, nil).Once()

		mockSocialAccountMySQLRepo.
			On("QueryById", mock.Anything, mock.MatchedBy(func(googleUUID string) bool { return googleUUID == mockGoogleUserInfo.ID })).
			Return(&mockSocialAccount, nil).Once()

		mockUserNySQLRepo.
			On("QueryById", mock.Anything, mock.MatchedBy(func(userID int64) bool { return userID == mockSocialAccount.UserID })).
			Return(&mockUser, nil).Once()


		s := NewSocialAccountUsecase(mockUserNySQLRepo, mockSocialAccountMySQLRepo, mockSocialAccountGoogleOAuth2Repo)
		_, err := s.GetUserInfo(context.TODO(), oauthToken)

		assert.NoError(t, err)
	})

	t.Run("Not First login Success", func(t *testing.T) {
		mockSocialAccountGoogleOAuth2Repo.
			On("GetUserInfo", mock.Anything, mock.MatchedBy(func(token *oauth2.Token) bool { return token == oauthToken })).
			Return(&mockGoogleUserInfo, nil).Once()

		mockSocialAccountMySQLRepo.
			On("QueryById", mock.Anything, mock.MatchedBy(func(googleUUID string) bool { return googleUUID == mockGoogleUserInfo.ID })).
			Return(nil, gorm.ErrRecordNotFound).Once()

		mockSocialAccountMySQLRepo.
			On("Store", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).Once()


		s := NewSocialAccountUsecase(mockUserNySQLRepo, mockSocialAccountMySQLRepo, mockSocialAccountGoogleOAuth2Repo)
		_, err := s.GetUserInfo(context.TODO(), oauthToken)

		assert.NoError(t, err)
	})

	t.Run("Get google user info Failed", func(t *testing.T) {
		mockSocialAccountGoogleOAuth2Repo.
			On("GetUserInfo", mock.Anything, mock.MatchedBy(func(token *oauth2.Token) bool { return token == oauthToken })).
			Return(nil, errors.New("get google user info failed")).Once()

		s := NewSocialAccountUsecase(mockUserNySQLRepo, mockSocialAccountMySQLRepo, mockSocialAccountGoogleOAuth2Repo)
		_, err := s.GetUserInfo(context.TODO(), oauthToken)

		assert.Equal(t, err, errors.New("get google user info failed"))
	})

	t.Run("Get social account Failed", func(t *testing.T) {
		mockSocialAccountGoogleOAuth2Repo.
			On("GetUserInfo", mock.Anything, mock.MatchedBy(func(token *oauth2.Token) bool { return token == oauthToken })).
			Return(&mockGoogleUserInfo, nil).Once()

		mockSocialAccountMySQLRepo.
			On("QueryById", mock.Anything, mock.MatchedBy(func(googleUUID string) bool { return googleUUID == mockGoogleUserInfo.ID })).
			Return(nil, errors.New("get social account info failed")).Once()

		s := NewSocialAccountUsecase(mockUserNySQLRepo, mockSocialAccountMySQLRepo, mockSocialAccountGoogleOAuth2Repo)
		_, err := s.GetUserInfo(context.TODO(), oauthToken)

		assert.Equal(t, err, errors.New("get social account info failed"))
	})
}
