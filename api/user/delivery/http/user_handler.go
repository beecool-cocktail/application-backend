package http

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)

type UserHandler struct {
	Configure   *service.Configure
	UserUsecase domain.UserUsecase
	SocialAccountUsecase domain.SocialAccountUsecase
}

func NewUserHandler(s *service.Service, userUsecase domain.UserUsecase, socialAccountUsecase domain.SocialAccountUsecase) {
	handler := &UserHandler{
		Configure:   s.Configure,
		UserUsecase: userUsecase,
		SocialAccountUsecase: socialAccountUsecase,
	}

	s.HTTP.GET("/api/google-login", handler.SocialLogin)
	s.HTTP.POST("/api/google-authenticate", handler.GoogleAuthenticate)
}

// swagger:route GET /google-login login googleLogin
//
// Login with google OAuth2
//
// Will redirect with authorization code.
//
// responses:
//  307: description: redirect
func (u *UserHandler) SocialLogin(c *gin.Context) {
	g := u.Configure.Others.GoogleOAuth2
	googleOAuth2Config := &oauth2.Config{
		ClientID:     g.ClientID,
		ClientSecret: g.ClientSecret,
		RedirectURL:  g.RedirectURL,
		Scopes:       g.Scopes,
		Endpoint:     google.Endpoint,
	}

	url := googleOAuth2Config.AuthCodeURL("whispering-corner")

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// swagger:operation POST /google-authenticate login googleAuthenticateRequest
// ---
// summary: Get access token.
// description: Use Code to exchange access token.
// responses:
//  "201":
//    "$ref": "#/responses/googleAuthenticateResponse"
func (u *UserHandler) GoogleAuthenticate(c *gin.Context) {
	var request viewmodels.GoogleAuthenticateRequest
	var response viewmodels.GoogleAuthenticateResponse

	if err := c.BindJSON(&request); err != nil {
		logrus.Error(err)
		util.PackResponseWithError(c, domain.ErrRequestDecodeFailed, "request unmarshal failed")
	}

	token, err := u.SocialAccountUsecase.Exchange(c, request.Code)
	if err != nil {
		util.PackResponseWithError(c, err, err.Error())
	}

	jwtToken, err := u.SocialAccountUsecase.GetUserInfo(c, token)
	if err != nil {
		util.PackResponseWithError(c, err, err.Error())
	}

	response.Token = jwtToken

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}