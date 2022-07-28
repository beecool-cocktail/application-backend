package http

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/cockarticletype"
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
	"strconv"
)

type UserHandler struct {
	Configure               *service.Configure
	Logger                  *logrus.Logger
	UserUsecase             domain.UserUsecase
	CocktailUsecase         domain.CocktailUsecase
	FavoriteCocktailUsecase domain.FavoriteCocktailUsecase
	SocialAccountUsecase    domain.SocialAccountUsecase
}

func NewUserHandler(s *service.Service, userUsecase domain.UserUsecase, socialAccountUsecase domain.SocialAccountUsecase,
	CocktailUsecase domain.CocktailUsecase, favoriteCocktailUsecase domain.FavoriteCocktailUsecase,
	middlewareHandler middleware.Handler) {

	handler := &UserHandler{
		Configure:               s.Configure,
		Logger:                  s.Logger,
		UserUsecase:             userUsecase,
		FavoriteCocktailUsecase: favoriteCocktailUsecase,
		CocktailUsecase:         CocktailUsecase,
		SocialAccountUsecase:    socialAccountUsecase,
	}

	s.HTTP.GET("/api/auth/google-login", handler.SocialLogin)
	s.HTTP.POST("/api/auth/google-authenticate", handler.GoogleAuthenticate)
	s.HTTP.POST("/api/auth/logout", handler.Logout)
	s.HTTP.GET("/api/users/current", middlewareHandler.JWTAuthMiddleware(), handler.GetUserInfo)
	s.HTTP.GET("/api/users/:userID", handler.GetOtherUserInfo)
	s.HTTP.PUT("/api/users/current", middlewareHandler.JWTAuthMiddleware(), handler.UpdateUserInfo)
	s.HTTP.POST("/api/users/current/favorite-cocktails", middlewareHandler.JWTAuthMiddleware(), handler.CollectArticle)
	s.HTTP.DELETE("/api/users/current/favorite-cocktails/:cocktailID", middlewareHandler.JWTAuthMiddleware(), handler.RemoveCollectionArticle)
	s.HTTP.GET("/api/users/current/favorite-cocktails", middlewareHandler.JWTAuthMiddleware(), handler.GetUserFavoriteList)
	s.HTTP.GET("/api/users/:userID/favorite-cocktails", middlewareHandler.JWTAuthMiddlewareIfExist(), handler.GetOtherUserFavoriteList)
	s.HTTP.GET("/api/users/current/cocktails", middlewareHandler.JWTAuthMiddleware(), handler.SelfCocktailList)
	s.HTTP.GET("/api/users/:userID/cocktails", middlewareHandler.JWTAuthMiddlewareIfExist(), handler.OtherCocktailList)
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

// swagger:operation POST /auth/logout login logoutRequest
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

// swagger:operation GET /users/current user info
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

	numberOfPost, err := u.CocktailUsecase.QueryFormalCountsByUserID(c, user.ID)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("query formal counts by user_id failed - %s",
			err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	numberOfCollection, err := u.FavoriteCocktailUsecase.QueryCountsByUserID(c, user.ID)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("query favotite counts by user_id failed - %s",
			err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	response = viewmodels.GetUserInfoResponse{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Photo:  user.Photo,
		Length: user.Length,
		Width:  user.Width,
		Coordinate: []viewmodels.Coordinate{
			{
				X: user.CoordinateX1,
				Y: user.CoordinateY1,
			},
			{
				X: user.CoordinateX2,
				Y: user.CoordinateY2,
			},
		},
		NumberOfPost:       numberOfPost,
		NumberOfCollection: numberOfCollection,
		IsCollectionPublic: user.IsCollectionPublic,
	}

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation GET /users/{id} user getOtherUserInfo
// ---
// summary: Get other user information.
// description: Get other user id, name, email, numberOfPost, numberOfCollection and photo.
//
// security:
// - Bearer: [apiKey]
//
// parameters:
// - name: id
//   in: path
//   required: true
//   type: integer
//   example: 1
//
// responses:
//  "200":
//    "$ref": "#/responses/getOtherUserInfoResponse"
func (u *UserHandler) GetOtherUserInfo(c *gin.Context) {
	api := "/user/info"
	var response viewmodels.GetOtherUserInfoResponse
	userID := c.Param("userID")

	userIDNumber, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	user, err := u.UserUsecase.QueryById(c, userIDNumber)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("query by id failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	numberOfPost, err := u.CocktailUsecase.QueryFormalCountsByUserID(c, user.ID)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("query formal counts by user_id failed - %s",
			err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	numberOfCollection, err := u.FavoriteCocktailUsecase.QueryCountsByUserID(c, user.ID)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("query favotite counts by user_id failed - %s",
			err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	response = viewmodels.GetOtherUserInfoResponse{
		UserID: user.ID,
		Name:   user.Name,
		Photo:  user.Photo,
		Length: user.Length,
		Width:  user.Width,
		Coordinate: []viewmodels.Coordinate{
			{
				X: user.CoordinateX1,
				Y: user.CoordinateY1,
			},
			{
				X: user.CoordinateX2,
				Y: user.CoordinateY2,
			},
		},
		NumberOfPost:       numberOfPost,
		NumberOfCollection: numberOfCollection,
		IsCollectionPublic: user.IsCollectionPublic,
	}

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation PUT /users/current user updateUserInfoRequest
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

	if len(request.Coordinate) != 2 {
		service.GetLoggerEntry(u.Logger, api, request).Errorf("parameter illegal")
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}
	err := u.UserUsecase.UpdateUserInfo(c,
		&domain.User{
			ID:                 userId,
			Name:               request.Name,
			IsCollectionPublic: request.IsCollectionPublic,
			Length:             request.Length,
			Width:              request.Width,
			CoordinateX1:       request.Coordinate[0].X,
			CoordinateY1:       request.Coordinate[0].Y,
			CoordinateX2:       request.Coordinate[1].X,
			CoordinateY2:       request.Coordinate[1].Y,
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

// swagger:operation POST /users/current/favorite-cocktails user collectArticleRequest
// ---
// summary: Add cocktail article to favorite list.
// description: Add cocktail article to favorite list.
//
// security:
// - Bearer: [apiKey]
//
// responses:
//  "201": success
func (u *UserHandler) CollectArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")
	api := "POST /users/favorite-cocktails"

	var request viewmodels.CollectArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		service.GetLoggerEntry(u.Logger, api, request).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	favoriteCocktail := domain.FavoriteCocktail{
		CocktailID: request.ID,
		UserID:     userId,
	}
	err := u.FavoriteCocktailUsecase.Store(c, &favoriteCocktail)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("collect article failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	util.PackResponseWithData(c, http.StatusCreated, nil, domain.GetErrorCode(nil), "")
}

// swagger:operation DELETE /users/current/favorite-cocktails/{id} user removeCollectionArticle
// ---
// summary: Remove cocktail article from favorite list.
// description: Remove cocktail article from favorite list.
//
// security:
// - Bearer: [apiKey]
//
// parameters:
// - name: id
//   in: path
//   required: true
//   type: integer
//   example: 123456
//
// responses:
//  "200":
//    "$ref": "#/responses/deleteFavoriteCocktailResponse"
func (u *UserHandler) RemoveCollectionArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")
	cocktailID := c.Param("cocktailID")
	api := "Delete /users/favorite-cocktails/" + cocktailID

	var response viewmodels.DeleteFavoriteCocktailResponse

	cocktailIDNumber, err := strconv.ParseInt(cocktailID, 10, 64)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	commandID, err := u.FavoriteCocktailUsecase.Delete(c, cocktailIDNumber, userId)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("delete article from favotite lost failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	response.CommandID = commandID

	util.PackResponseWithData(c, http.StatusCreated, response, domain.GetErrorCode(nil), "")
}

// swagger:operation GET /users/current/favorite-cocktails user getUserFavoriteList
// ---
// summary: Get current user favorite cocktail article list.
// description: Get current user favorite cocktail article list.
//
// security:
// - Bearer: [apiKey]
//
// responses:
//  "200":
//    "$ref": "#/responses/getUserFavoriteCocktailListResponse"
func (u *UserHandler) GetUserFavoriteList(c *gin.Context) {
	api := "/users/current/favorite-cocktails"
	var response viewmodels.GetUserFavoriteCocktailListResponse
	userId := c.GetInt64("user_id")

	favoriteCocktails, total, err := u.FavoriteCocktailUsecase.QueryByUserID(c, userId, domain.PaginationUsecase{},
		0)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("query by id failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	list := make([]viewmodels.FavoriteCocktail, 0)
	for _, cocktail := range favoriteCocktails {
		out := viewmodels.FavoriteCocktail{
			CocktailID: cocktail.CocktailID,
			UserName:   cocktail.UserName,
			Photo:      cocktail.CoverPhoto,
			Title:      cocktail.Title,
		}

		list = append(list, out)
	}

	response = viewmodels.GetUserFavoriteCocktailListResponse{
		Total:                total,
		FavoriteCocktailList: list,
	}

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation GET /users/{id}/favorite-cocktails user getOtherUserFavoriteList
// ---
// summary: Get other user favorite cocktail article list.
// description: Get other user favorite cocktail article list.
//
// security:
// - Bearer: [apiKey]
//
// parameters:
// - name: id
//   in: path
//   required: true
//   type: integer
//   example: 123456
//
// responses:
//  "200":
//    "$ref": "#/responses/getUserFavoriteCocktailListResponse"
func (u *UserHandler) GetOtherUserFavoriteList(c *gin.Context) {
	api := "/users/{id}/favorite-cocktails"
	userID := c.Param("userID")
	selfUserID := c.GetInt64("user_id")

	var response viewmodels.GetUserFavoriteCocktailListResponse

	userIDNumber, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	user, err := u.UserUsecase.QueryById(c, userIDNumber)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("query user by user id failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	list := make([]viewmodels.FavoriteCocktail, 0)
	if !user.IsCollectionPublic {
		response = viewmodels.GetUserFavoriteCocktailListResponse{
			IsPublic:             user.IsCollectionPublic,
			Total:                0,
			FavoriteCocktailList: list,
		}
		util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
		return
	}

	favoriteCocktails, total, err := u.FavoriteCocktailUsecase.QueryByUserID(c, userIDNumber,
		domain.PaginationUsecase{}, selfUserID)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("query favorite cocktail by user id failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	for _, cocktail := range favoriteCocktails {
		out := viewmodels.FavoriteCocktail{
			CocktailID:  cocktail.CocktailID,
			UserName:    cocktail.UserName,
			Photo:       cocktail.CoverPhoto,
			Title:       cocktail.Title,
			IsCollected: cocktail.IsCollected,
		}

		list = append(list, out)
	}

	response = viewmodels.GetUserFavoriteCocktailListResponse{
		IsPublic:             user.IsCollectionPublic,
		Total:                total,
		FavoriteCocktailList: list,
	}

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation GET /users/current/cocktails user selfCocktailList
// ---
// summary: Get self cocktail list
// description: Get self cocktail list order by create date.
//
// security:
// - Bearer: [apiKey]
//
// responses:
//  "200":
//    "$ref": "#/responses/getSelfCocktailListResponse"
func (u *UserHandler) SelfCocktailList(c *gin.Context) {
	api := "/self-cocktails"
	userId := c.GetInt64("user_id")

	var response viewmodels.GetSelfCocktailListResponse

	filter := make(map[string]interface{})
	filter["category"] = cockarticletype.Formal
	filter["user_id"] = userId
	cocktails, err := u.CocktailUsecase.QueryFormalByUserID(c, userId, 0)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("get cocktails with filter failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	cocktailList := make([]viewmodels.SelfCocktailList, 0)
	for _, cocktail := range cocktails {
		ingredients := make([]viewmodels.CocktailIngredientWithoutIDInResponse, 0)
		for _, ingredient := range cocktail.Ingredients {
			out := viewmodels.CocktailIngredientWithoutIDInResponse{
				Name:   ingredient.IngredientName,
				Amount: ingredient.IngredientAmount,
			}
			ingredients = append(ingredients, out)
		}

		out := viewmodels.SelfCocktailList{
			CocktailID: cocktail.CocktailID,
			UserName:   cocktail.UserName,
			Title:      cocktail.Title,
			Photo:      cocktail.CoverPhoto.Photo,
		}

		cocktailList = append(cocktailList, out)
	}

	response.CocktailList = cocktailList

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation GET /users/{id}/cocktails user otherCocktailList
// ---
// summary: Get other user cocktail list
// description: Get other user cocktail list order by create date.
//
// security:
// - Bearer: [apiKey]
//
// parameters:
// - name: id
//   in: path
//   required: true
//   type: integer
//   example: 1
//
// responses:
//  "200":
//    "$ref": "#/responses/getOtherCocktailListResponse"
func (u *UserHandler) OtherCocktailList(c *gin.Context) {
	api := "/self-cocktails"
	targetUserID := c.Param("userID")
	selfUserID := c.GetInt64("user_id")

	userIDNumber, err := strconv.ParseInt(targetUserID, 10, 64)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	var response viewmodels.GetOtherCocktailListResponse

	filter := make(map[string]interface{})
	filter["category"] = cockarticletype.Formal
	filter["user_id"] = userIDNumber
	cocktails, err := u.CocktailUsecase.QueryFormalByUserID(c, userIDNumber, selfUserID)
	if err != nil {
		service.GetLoggerEntry(u.Logger, api, nil).Errorf("get cocktails with filter failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	cocktailList := make([]viewmodels.OtherCocktailList, 0)
	for _, cocktail := range cocktails {
		ingredients := make([]viewmodels.CocktailIngredientWithoutIDInResponse, 0)
		for _, ingredient := range cocktail.Ingredients {
			out := viewmodels.CocktailIngredientWithoutIDInResponse{
				Name:   ingredient.IngredientName,
				Amount: ingredient.IngredientAmount,
			}
			ingredients = append(ingredients, out)
		}

		out := viewmodels.OtherCocktailList{
			CocktailID:  cocktail.CocktailID,
			UserName:    cocktail.UserName,
			Title:       cocktail.Title,
			IsCollected: cocktail.IsCollected,
			Photo:       cocktail.CoverPhoto.Photo,
		}

		cocktailList = append(cocktailList, out)
	}

	response.CocktailList = cocktailList

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}
