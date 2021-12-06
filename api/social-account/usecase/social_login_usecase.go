package usecase

import (
	"context"
	"errors"
	"github.com/beecool-cocktail/application-backend/domain"
	social_account "github.com/beecool-cocktail/application-backend/enum/social-account"
	"github.com/beecool-cocktail/application-backend/middleware"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type socialLoginUsecase struct {
	userMySQLRepo                domain.UserMySQLRepository
	socialAccountMySQLRepo       domain.SocialAccountMySQLRepository
	socialAccountGoogleOAuthRepo domain.SocialAccountGoogleOAuthRepository
}

func NewSocialAccountUsecase(userMySQLRepo domain.UserMySQLRepository, socialAccountMySQLRepo domain.SocialAccountMySQLRepository,
	socialAccountGoogleOAuthRepo domain.SocialAccountGoogleOAuthRepository) domain.SocialAccountUsecase {
	return &socialLoginUsecase{
		userMySQLRepo:                userMySQLRepo,
		socialAccountMySQLRepo:       socialAccountMySQLRepo,
		socialAccountGoogleOAuthRepo: socialAccountGoogleOAuthRepo,
	}
}

func (s *socialLoginUsecase) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.socialAccountGoogleOAuthRepo.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *socialLoginUsecase) GetUserInfo(ctx context.Context, token *oauth2.Token) (string, error) {
	googleUserInfo, err := s.socialAccountGoogleOAuthRepo.GetUserInfo(ctx, token)
	if err != nil {
		return "", err
	}

	var userID int64
	var account, name string
	socialAccount, err := s.socialAccountMySQLRepo.QueryById(ctx, googleUserInfo.ID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Error(err)
			return "", err
		} else {
			// first time use google login, register a new user
			account = util.GenString(16)
			userID = util.GetID(util.IdGenerator)
			name = googleUserInfo.Name
			storeErr := s.socialAccountMySQLRepo.Store(ctx,
				&domain.SocialAccount{
					SocialID: googleUserInfo.ID,
					UserID:   userID,
					Type:     social_account.ParseSocialAccountType(social_account.Google),
				},
				&domain.User{
					UserID:  userID,
					Account: account,
					Name:    name,
					Email:   googleUserInfo.Email,
				})
			if storeErr != nil {
				return "", err
			}
		}
	} else {
		// not first time use google login, todo something
		user, err := s.userMySQLRepo.QueryById(ctx, socialAccount.UserID)
		if err != nil {
			return "", nil
		}

		account = user.Account
		userID = user.UserID
		name = user.Name
	}

	payloadData := middleware.PayloadData{
		UserID:  userID,
		Account: account,
		Name:    name,
	}
	jwtToken, err := middleware.GenToken(payloadData)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
