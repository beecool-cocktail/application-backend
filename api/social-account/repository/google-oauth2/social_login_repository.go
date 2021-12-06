package google_oauth2

import (
	"context"
	"encoding/json"
	"github.com/beecool-cocktail/application-backend/domain"
	"golang.org/x/oauth2"
	"io/ioutil"
)

type socialAccountGoogleOAuthRepository struct {
	googleOAuthConfig *oauth2.Config
}

func NewGoogleOAuthSocialAccountRepository(googleOAuthConfig *oauth2.Config) domain.SocialAccountGoogleOAuthRepository {
	return &socialAccountGoogleOAuthRepository{googleOAuthConfig}
}

func (s *socialAccountGoogleOAuthRepository) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.googleOAuthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}

	return token, err
}

func (s *socialAccountGoogleOAuthRepository) GetUserInfo(ctx context.Context, token *oauth2.Token) (*domain.GoogleUserInfo, error) {
	client := s.googleOAuthConfig.Client(oauth2.NoContext, token)

	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var userInfo domain.GoogleUserInfo
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}