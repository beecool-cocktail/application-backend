package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/beecool-cocktail/application-backend/domain"
	social_account "github.com/beecool-cocktail/application-backend/enum/social-account"
	"github.com/beecool-cocktail/application-backend/enum/usertype"
	"github.com/beecool-cocktail/application-backend/middleware"
	"github.com/beecool-cocktail/application-backend/util"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type socialLoginUsecase struct {
	userMySQLRepo                domain.UserMySQLRepository
	userRedisRepo                domain.UserRedisRepository
	socialAccountMySQLRepo       domain.SocialAccountMySQLRepository
	socialAccountRedisRepo 		 domain.SocialAccountRedisRepository
	socialAccountGoogleOAuthRepo domain.SocialAccountGoogleOAuthRepository
}

func NewSocialAccountUsecase(userMySQLRepo domain.UserMySQLRepository, userRedisRepo domain.UserRedisRepository,
	socialAccountMySQLRepo domain.SocialAccountMySQLRepository, socialAccountRedisRepo 		 domain.SocialAccountRedisRepository, 
	socialAccountGoogleOAuthRepo domain.SocialAccountGoogleOAuthRepository) domain.SocialAccountUsecase {
	return &socialLoginUsecase{
		userMySQLRepo:                userMySQLRepo,
		userRedisRepo:                userRedisRepo,
		socialAccountMySQLRepo:       socialAccountMySQLRepo,
		socialAccountRedisRepo: 		 socialAccountRedisRepo,
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
	var userType int
	var account, name string
	socialAccount, err := s.socialAccountMySQLRepo.QueryById(ctx, googleUserInfo.Sub)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		} else {
			// first time use google login, register a new user
			account = util.GenString(16)
			name = googleUserInfo.Name
			userType = usertype.Normal.Int()
			userID, err = s.socialAccountMySQLRepo.Store(ctx,
				&domain.SocialAccount{
					SocialID: googleUserInfo.Sub,
					Type:     social_account.ParseSocialAccountType(social_account.Google),
				},
				&domain.User{
					Account: account,
					Name:    name,
					Email:   googleUserInfo.Email,
				})
			if err != nil {
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
		userID = user.ID
		name = user.Name
		userType = user.Type
	}

	payloadData := middleware.PayloadData{
		UserID:  userID,
		Account: account,
		Name:    name,
		Type:    userType,
	}
	jwtToken, err := middleware.GenToken(payloadData)
	if err != nil {
		return "", err
	}

	err = s.userRedisRepo.Store(ctx, &domain.UserCache{
		Id:          userID,
		Account:     account,
		Name:        name,
		AccessToken: jwtToken,
	})
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (s *socialLoginUsecase) GenerateState(ctx context.Context, state domain.State) (string, error) {
	stateString, err := getState(32)
	if err != nil {
		return "", err
	}

	err = s.socialAccountRedisRepo.StoreState(ctx, stateString, state)
	if err != nil {
		return "", err
	}

	return stateString, nil
}

func getState(n int) (string, error) {
	data := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (s *socialLoginUsecase) GetState(ctx context.Context, stateString string) (domain.State, error) {
	state, err := s.socialAccountRedisRepo.GetState(ctx, stateString)
	if err != nil {
		return domain.State{}, err
	}

	return state, nil
}