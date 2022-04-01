package http

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/middleware"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vincent-petithory/dataurl"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)

type UserHandler struct {
	Configure            *service.Configure
	Logger               *logrus.Logger
	UserUsecase          domain.UserUsecase
	SocialAccountUsecase domain.SocialAccountUsecase
}

func NewUserHandler(s *service.Service, userUsecase domain.UserUsecase, socialAccountUsecase domain.SocialAccountUsecase,
	middlewareHandler middleware.Handler) {

	handler := &UserHandler{
		Configure:            s.Configure,
		Logger:               s.Logger,
		UserUsecase:          userUsecase,
		SocialAccountUsecase: socialAccountUsecase,
	}

	s.HTTP.GET("/api/auth/google-login", handler.SocialLogin)
	s.HTTP.POST("/api/auth/google-authenticate", handler.GoogleAuthenticate)
	s.HTTP.POST("/api/auth/logout", handler.Logout)
	s.HTTP.GET("/api/user/current", middlewareHandler.JWTAuthMiddleware(), handler.GetUserInfo)
	s.HTTP.PUT("/api/user/current", middlewareHandler.JWTAuthMiddleware(), handler.UpdateUserInfo)
}

// swagger:route GET /auth/google-login login googleLogin
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

// swagger:operation POST /auth/google-authenticate login googleAuthenticateRequest
// ---
// summary: Get access token.
// description: Use Code to exchange access token.
// responses:
//  "201":
//    "$ref": "#/responses/googleAuthenticateResponse"
func (u *UserHandler) GoogleAuthenticate(c *gin.Context) {
	api := "/google-authenticate"
	var request viewmodels.GoogleAuthenticateRequest
	var response viewmodels.GoogleAuthenticateResponse

	if err := c.BindJSON(&request); err != nil {
		service.GetLoggerEntry(u.Logger, api, request).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	service.GetLoggerEntry(u.Logger, api, request).Info()

	token, err := u.SocialAccountUsecase.Exchange(c, request.Code)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, request).Errorf("exchange google token failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	jwtToken, err := u.SocialAccountUsecase.GetUserInfo(c, token)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, request).Errorf("get user info failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	response.Token = jwtToken
	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation POST /auth/logout user logoutRequest
// ---
// summary: User logout.
// description: make token invalid.
// responses:
//  "200": success
func (u *UserHandler) Logout(c *gin.Context) {
	api := "/user/logout"
	var request viewmodels.LogoutRequest

	if err := c.BindJSON(&request); err != nil {
		service.GetLoggerEntry(u.Logger, api, request).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	err := u.UserUsecase.Logout(c, request.UserID)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, request).Errorf("logout failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	util.PackResponseWithData(c, http.StatusOK, nil, domain.GetErrorCode(nil), "")
}

// swagger:operation GET /user/current user info
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
	api := "/user/info"
	var response viewmodels.GetUserInfoResponse
	userId := c.GetInt64("user_id")

	user, err := u.UserUsecase.QueryById(c, userId)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("query by id failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	response = viewmodels.GetUserInfoResponse{
		UserID:             user.ID,
		Name:               user.Name,
		Email:              user.Email,
		Photo:              user.Photo,
		NumberOfPost:       user.NumberOfPost,
		NumberOfCollection: user.NumberOfCollection,
		IsCollectionPublic: user.IsCollectionPublic,
	}

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation PUT /user/current user updateUserInfoRequest
// ---
// summary: Edit user information.
// description: Edit user name and collection of publicity status.
//
// security:
// - Bearer: [apiKey]
//
// responses:
//  "200":
//    "$ref": "#/responses/updateUserPhotoResponse"
func (u *UserHandler) UpdateUserInfo(c *gin.Context) {
	api := "/user/edit-info"

	var request viewmodels.UpdateUserInfoRequest
	var response viewmodels.UpdateUserInfoResponse
	var userImage domain.UserImage
	if err := c.ShouldBindJSON(&request); err != nil {
		service.GetLoggerEntry(u.Logger, api, request).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	userId := c.GetInt64("user_id")

	//user update photo
	if request.File != "" {
		dataURL, err := dataurl.DecodeString(request.File)
		if err != nil {
			service.GetLoggerEntry(u.Logger, api, request).Errorf("decode data url failed - %s", err)
			util.PackResponseWithError(c, err, err.Error())
			return
		}
		userImage = domain.UserImage{
			ID:   userId,
			Data: string(dataURL.Data),
			Type: dataURL.MediaType.ContentType(),
		}
	} else {
		// user didn't update photo
	}

	err := u.UserUsecase.UpdateUserInfo(c,
		&domain.User{
			ID:                 userId,
			Name:               request.Name,
			IsCollectionPublic: request.IsCollectionPublic,
		},
		&userImage)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("update basic info failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	response = viewmodels.UpdateUserInfoResponse{
		Photo: userImage.Destination,
	}

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}
