package http

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SocialAccount struct {
	SocialAccountUsecase domain.SocialAccountUsecase
}

func NewSocialAccountHandler(s *service.Service, socialAccountUsecase domain.SocialAccountUsecase) {
	handler := &SocialAccount{
		SocialAccountUsecase: socialAccountUsecase,
	}

	s.HTTP.POST("/api/google-authenticate", handler.GoogleAuthenticate)
}

func (s *SocialAccount) GoogleAuthenticate(c *gin.Context) {
	state := c.Request.FormValue("state")
	if state != "whispering-corner" {

	}
	code := c.Request.FormValue("code")
	token, err := s.SocialAccountUsecase.Exchange(c, code)
	if err != nil {
		util.PackResponseWithError(c, err, err.Error())
	}

	jwtToken, err := s.SocialAccountUsecase.GetUserInfo(c, token)
	if err != nil {
		util.PackResponseWithError(c, err, err.Error())
	}

	c.Header("Authorization", jwtToken)
	c.Redirect(http.StatusTemporaryRedirect, "localhost:3000")
}