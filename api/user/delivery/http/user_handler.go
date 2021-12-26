package http

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/middleware"
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

func NewUserHandler(s *service.Service, userUsecase domain.UserUsecase, socialAccountUsecase domain.SocialAccountUsecase,
	middlewareHandler middleware.Handler) {

	handler := &UserHandler{
		Configure:   s.Configure,
		UserUsecase: userUsecase,
		SocialAccountUsecase: socialAccountUsecase,
	}

	s.HTTP.GET("/api/google-login", handler.SocialLogin)
	s.HTTP.POST("/api/google-authenticate", handler.GoogleAuthenticate)
	s.HTTP.POST("/api/user/logout", handler.Logout)
	s.HTTP.GET("/api/user/info", middlewareHandler.JWTAuthMiddleware(), handler.GetUserInfo)
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
		return
	}

	token, err := u.SocialAccountUsecase.Exchange(c, request.Code)
	if err != nil {
		logrus.Error(err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	jwtToken, err := u.SocialAccountUsecase.GetUserInfo(c, token)
	if err != nil {
		logrus.Error(err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	response.Token = jwtToken
	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation POST /user/logout user logoutRequest
// ---
// summary: User logout.
// description: make token invalid.
// responses:
//  "200": success
func (u *UserHandler) Logout(c *gin.Context) {
	var request viewmodels.LogoutRequest

	if err := c.BindJSON(&request); err != nil {
		logrus.Error(err)
		util.PackResponseWithError(c, domain.ErrRequestDecodeFailed, "request unmarshal failed")
		return
	}

	err := u.UserUsecase.Logout(c, request.UserID)
	if err != nil {
		util.PackResponseWithError(c, err, err.Error())
		return
	}


	util.PackResponseWithData(c, http.StatusOK, nil, domain.GetErrorCode(nil), "")
}

// swagger:operation GET /user/info user info
// ---
// summary: Get user information.
// description: Get user id, name, email, numberOfPost, numberOfCollection and photo.
//
// security:
// - Bearer: [apiKey]
//
// responses:
//  "200":
//    "$ref": "#/responses/getUserInfoResponse"
func (u *UserHandler) GetUserInfo(c *gin.Context) {
	var response viewmodels.GetUserInfoResponse
	userId := c.GetInt64("user_id")

	user, err := u.UserUsecase.QueryById(c, userId)
	if err != nil {
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	response = viewmodels.GetUserInfoResponse{
		UserID: user.ID,
		Name: user.Name,
		Email: user.Email,
		Photo: user.Photo,
		NumberOfPost: user.NumberOfPost,
		NumberOfCollection: user.NumberOfCollection,
	}

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}