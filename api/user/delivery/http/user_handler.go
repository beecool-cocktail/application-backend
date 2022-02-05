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
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
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

	s.HTTP.GET("/api/google-login", handler.SocialLogin)
	s.HTTP.POST("/api/google-authenticate", handler.GoogleAuthenticate)
	s.HTTP.POST("/api/user/logout", handler.Logout)
	s.HTTP.GET("/api/user/info", middlewareHandler.JWTAuthMiddleware(), handler.GetUserInfo)
	s.HTTP.POST("/api/user/edit-info", middlewareHandler.JWTAuthMiddleware(), handler.UpdateUserInfo)
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

// swagger:operation POST /user/logout user logoutRequest
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

// swagger:operation POST /user/edit-info user updateUserInfoRequest
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
	var response viewmodels.UpdateUserPhotoResponse
	var userImage domain.UserImage

	userId := c.GetInt64("user_id")

	form, err := c.MultipartForm()
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, form).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	var file *multipart.FileHeader

	//Todo move to another function
	files := form.File
	if _, ok := files["file"]; ok {
		//user update photo
		if len(files["file"]) > 0 {
			file = files["file"][0]
			userImage = domain.UserImage{
				ID:   userId,
				Data: file,
				Type: filepath.Ext(file.Filename),
			}
		} else {
			service.GetLoggerEntry(u.Logger, api, form).Errorf("parameter illegal - %s", err)
			util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
			return
		}
	} else {
		//user didn't update photo
	}

	values := form.Value
	isCollectionPublic, err := strconv.ParseBool(values["is_collection_public"][0])
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, form).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	if _, ok := values["name"]; !ok && len(values["name"]) <= 0 {
		service.GetLoggerEntry(u.Logger, api, form).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	name := values["name"][0]

	err = u.UserUsecase.UpdateUserInfo(c,
		&domain.User{
		ID:                 userId,
		Name:               name,
		IsCollectionPublic: isCollectionPublic,
		},
		&userImage)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("update basic info failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	response = viewmodels.UpdateUserPhotoResponse{
		Photo: userImage.Destination,
	}

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}
