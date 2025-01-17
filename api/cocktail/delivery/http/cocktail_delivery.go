package http

import (
	"errors"
	"net/http"

	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/cockarticletype"
	"github.com/beecool-cocktail/application-backend/enum/usertype"
	"github.com/beecool-cocktail/application-backend/middleware"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vincent-petithory/dataurl"
)

type CocktailHandler struct {
	Service         *service.Service
	CocktailUsecase domain.CocktailUsecase
	UserUsecase     domain.UserUsecase
}

func NewCocktailHandler(s *service.Service, cocktailUsecase domain.CocktailUsecase,
	userUsecase domain.UserUsecase,
	middlewareHandler middleware.Handler) {
	handler := &CocktailHandler{
		Service:         s,
		UserUsecase:     userUsecase,
		CocktailUsecase: cocktailUsecase,
	}

	s.HTTP.GET("/api/cocktails/:id", middlewareHandler.JWTAuthMiddlewareIfExist(), handler.GetCocktailByCocktailID)
	s.HTTP.GET("/api/cocktail-drafts/:id", middlewareHandler.JWTAuthMiddleware(), handler.GetCocktailDraftByCocktailID)
	s.HTTP.GET("/api/cocktails", middlewareHandler.JWTAuthMiddlewareIfExist(), handler.CocktailList)
	s.HTTP.GET("/api/cocktail-drafts", middlewareHandler.JWTAuthMiddleware(), handler.CocktailDraftList)
	s.HTTP.POST("/api/cocktails", middlewareHandler.JWTAuthMiddleware(), handler.PostArticle)
	s.HTTP.POST("/api/cocktail-drafts", middlewareHandler.JWTAuthMiddleware(), handler.PostDraftArticle)
	s.HTTP.POST("/api/cocktail-drafts/:id", middlewareHandler.JWTAuthMiddleware(), handler.MakeDraftArticleToFormalArticle)
	s.HTTP.PUT("/api/cocktails/:id", middlewareHandler.JWTAuthMiddleware(), handler.UpdateFormalArticle)
	s.HTTP.PUT("/api/cocktail-drafts/:id", middlewareHandler.JWTAuthMiddleware(), handler.UpdateDraftArticle)
	s.HTTP.PATCH("/api/cocktail-drafts", middlewareHandler.JWTAuthMiddleware(), handler.DeleteDraftArticle)
	s.HTTP.PATCH("/api/cocktails", middlewareHandler.JWTAuthMiddleware(), handler.DeleteFormalArticle)
}

// swagger:operation GET /cocktails/{id} cocktail getCocktailByIDRequest
// ---
// summary: Get cocktail detail information
// description: Get cocktail photo, steps and ingredient.
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
//    "$ref": "#/responses/getCocktailByIDResponse"

func (co *CocktailHandler) GetCocktailByCocktailID(c *gin.Context) {
	var request viewmodels.GetCocktailByIDRequest
	var response viewmodels.GetCocktailByIDResponse
	userId := c.GetInt64("user_id")
	loggerFields := co.Service.Logger.GetLoggerFields(userId, c.ClientIP(), c.Request.Method, nil,
		c.Request.RequestURI)

	err := c.ShouldBindUri(&request)
	if err != nil {
		co.Service.Logger.LogFile(c, logrus.InfoLevel, loggerFields, "parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
	}

	cocktail, err := co.CocktailUsecase.QueryByCocktailID(c, request.ID, userId)
	if err != nil {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields, "query by cocktail id failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	ingredients := make([]viewmodels.CocktailIngredientWithoutIDInResponse, 0)
	for _, ingredient := range cocktail.Ingredients {
		out := viewmodels.CocktailIngredientWithoutIDInResponse{
			Name:   ingredient.IngredientName,
			Amount: ingredient.IngredientAmount,
		}
		ingredients = append(ingredients, out)
	}

	steps := make([]viewmodels.CocktailStepWithoutIDInResponse, 0)
	for _, step := range cocktail.Steps {
		out := viewmodels.CocktailStepWithoutIDInResponse{
			Description: step.StepDescription,
		}
		steps = append(steps, out)
	}

	lowQualityPhotos := make([]string, 0)
	for _, photo := range cocktail.LowQualityPhotos {
		fileName, err := util.GetFileNameByPath(photo.Photo)
		if err != nil {
			co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
				"get file name by path failed - %s", err)
			util.PackResponseWithError(c, err, err.Error())
			return
		}
		pathInServer := util.ConcatString("/", co.Service.Configure.Others.File.Image.PathInServer, fileName)
		dataURL, err := util.ParseLQIPFileToDataURL(pathInServer)
		if err != nil {
			co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
				"parse lqip file to data url failed - %s", err)
			util.PackResponseWithError(c, err, err.Error())
			return
		}
		lowQualityPhotos = append(lowQualityPhotos, dataURL)
	}

	photos := make([]viewmodels.CocktailPhotoWithIDInResponse, 0)
	for i, photo := range cocktail.Photos {
		out := viewmodels.CocktailPhotoWithIDInResponse{
			ID:               photo.ID,
			ImagePath:        photo.Photo,
			BlurImageDataURL: lowQualityPhotos[i],
		}
		photos = append(photos, out)
	}

	cocktailUser, err := co.UserUsecase.QueryById(c, cocktail.UserID)
	if err != nil {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
			"query user by user id failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	response = viewmodels.GetCocktailByIDResponse{
		CocktailID: cocktail.CocktailID,
		UserID:     cocktail.UserID,
		UserName:   cocktail.UserName,
		UserPhoto:  cocktailUser.CropAvatar,
		Height:     cocktailUser.Height,
		Width:      cocktailUser.Width,
		Coordinate: []viewmodels.Coordinate{
			{
				X: cocktailUser.CoordinateX1,
				Y: cocktailUser.CoordinateY1,
			},
			{
				X: cocktailUser.CoordinateX2,
				Y: cocktailUser.CoordinateY2,
			},
		},
		Rotation:       cocktailUser.Rotation,
		Title:          cocktail.Title,
		Description:    cocktail.Description,
		IngredientList: ingredients,
		StepList:       steps,
		Photos:         photos,
		IsCollected:    cocktail.IsCollected,
		CreatedDate:    cocktail.CreatedDate,
	}

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation GET /cocktail-drafts/{id} cocktail getCocktailDraftByIDRequest
// ---
// summary: Get cocktail draft detail information
// description: Get cocktail draft photo, steps and ingredient.
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
//    "$ref": "#/responses/getCocktailDraftByIDResponse"

func (co *CocktailHandler) GetCocktailDraftByCocktailID(c *gin.Context) {
	var request viewmodels.GetCocktailDraftByIDRequest
	var response viewmodels.GetCocktailDraftByIDResponse
	userId := c.GetInt64("user_id")
	loggerFields := co.Service.Logger.GetLoggerFields(userId, c.ClientIP(), c.Request.Method, nil,
		c.Request.RequestURI)

	err := c.ShouldBindUri(&request)
	if err != nil {
		co.Service.Logger.LogFile(c, logrus.InfoLevel, loggerFields, "parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
	}
	cocktail, err := co.CocktailUsecase.QueryDraftByCocktailID(c, request.ID, userId)
	if err != nil {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields, "query by cocktail id failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	ingredients := make([]viewmodels.CocktailIngredientWithoutIDInResponse, 0)
	for _, ingredient := range cocktail.Ingredients {
		out := viewmodels.CocktailIngredientWithoutIDInResponse{
			Name:   ingredient.IngredientName,
			Amount: ingredient.IngredientAmount,
		}
		ingredients = append(ingredients, out)
	}

	steps := make([]viewmodels.CocktailStepWithoutIDInResponse, 0)
	for _, step := range cocktail.Steps {
		out := viewmodels.CocktailStepWithoutIDInResponse{
			Description: step.StepDescription,
		}
		steps = append(steps, out)
	}

	photos := make([]viewmodels.DraftCocktailPhotoWithIDInResponse, 0)
	for _, photo := range cocktail.Photos {
		out := viewmodels.DraftCocktailPhotoWithIDInResponse{
			ID:        photo.ID,
			ImagePath: photo.Photo,
		}
		photos = append(photos, out)
	}

	response = viewmodels.GetCocktailDraftByIDResponse{
		CocktailID:     cocktail.CocktailID,
		Title:          cocktail.Title,
		Description:    cocktail.Description,
		IngredientList: ingredients,
		StepList:       steps,
		Photos:         photos,
		CreatedDate:    cocktail.CreatedDate,
	}

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation GET /cocktails cocktail getCocktail
// ---
// summary: Get popular cocktail list
// description: Get popular cocktail list order by create date.
//
// parameters:
// - name: page
//   in: query
//   required: true
//   type: integer
//   example: 1
//
// - name: page_size
//   in: query
//   required: true
//   type: integer
//   example: 10
//
// - name: keyword
//   in: query
//   required: false
//   type: string
//   example: search
//
// responses:
//  "200":
//	   "$ref": "#/responses/popularCocktailListResponse"

func (co *CocktailHandler) CocktailList(c *gin.Context) {
	userId := c.GetInt64("user_id")

	var request viewmodels.GetPopularCocktailListRequest
	var response viewmodels.GetPopularCocktailListResponse

	err := c.ShouldBind(&request)
	if err != nil {
		co.Service.Logger.LogFile(c, logrus.InfoLevel, co.Service.Logger.GetLoggerFields(userId, c.ClientIP(),
			c.Request.Method, nil, c.Request.RequestURI), "parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
	}
	loggerFields := co.Service.Logger.GetLoggerFields(userId, c.ClientIP(), c.Request.Method, request,
		c.Request.RequestURI)

	var cocktails []domain.APICocktail
	var total int64
	if request.Keyword != "" && co.Service.Configure.Elastic.Enable {
		cocktails, total, err = co.CocktailUsecase.Search(c, request.Keyword, request.Page, request.PageSize, userId)
		if err != nil {
			co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
				"get cocktails with keyword failed - %s", err)
			util.PackResponseWithError(c, err, err.Error())
			return
		}
	} else {
		filter := make(map[string]interface{})
		filter["category"] = cockarticletype.Formal
		cocktails, total, err = co.CocktailUsecase.GetAllWithFilter(c, filter, domain.PaginationUsecase{
			Page:     request.Page,
			PageSize: request.PageSize,
		},
			userId)
		if err != nil {
			co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
				"get cocktails with filter failed - %s", err)
			util.PackResponseWithError(c, err, err.Error())
			return
		}
	}

	cocktailList := make([]viewmodels.PopularCocktailList, 0)
	for _, cocktail := range cocktails {
		ingredients := make([]viewmodels.CocktailIngredientWithoutIDInResponse, 0)
		for _, ingredient := range cocktail.Ingredients {
			out := viewmodels.CocktailIngredientWithoutIDInResponse{
				Name:   ingredient.IngredientName,
				Amount: ingredient.IngredientAmount,
			}
			ingredients = append(ingredients, out)
		}

		lowQualityPhotos := make([]string, 0)
		for _, photo := range cocktail.LowQualityPhotos {
			fileName, err := util.GetFileNameByPath(photo.Photo)
			if err != nil {
				co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
					"get file name by path failed - %s", err)
				util.PackResponseWithError(c, err, err.Error())
				return
			}

			pathInServer := util.ConcatString("/", co.Service.Configure.Others.File.Image.PathInServer, fileName)
			dataURL, err := util.ParseLQIPFileToDataURL(pathInServer)
			if err != nil {
				co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
					"parse lqip file to data url failed - %s", err)
				util.PackResponseWithError(c, err, err.Error())
				return
			}
			lowQualityPhotos = append(lowQualityPhotos, dataURL)
		}

		photos := make([]viewmodels.CocktailPhotoWithIDInResponse, 0)
		for i, photo := range cocktail.Photos {
			out := viewmodels.CocktailPhotoWithIDInResponse{
				ID:               cocktail.CocktailID,
				ImagePath:        photo.Photo,
				BlurImageDataURL: lowQualityPhotos[i],
			}
			photos = append(photos, out)
		}

		out := viewmodels.PopularCocktailList{
			CocktailID:     cocktail.CocktailID,
			UserID:         cocktail.UserID,
			UserName:       cocktail.UserName,
			Title:          cocktail.Title,
			Photos:         photos,
			IngredientList: ingredients,
			IsCollected:    cocktail.IsCollected,
			CreatedDate:    cocktail.CreatedDate,
		}

		cocktailList = append(cocktailList, out)
	}

	response.Total = total
	response.PopularCocktailList = cocktailList

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation GET /cocktail-drafts cocktail getCocktailDraft
// ---
// summary: Get draft cocktail list
// description: Get draft cocktail list order by create date.
//
// security:
// - Bearer: [apiKey]
//
// parameters:
//
// responses:
//   "200":
//     "$ref": "#/responses/getDraftCocktailListResponse"

func (co *CocktailHandler) CocktailDraftList(c *gin.Context) {
	var response viewmodels.GetDraftCocktailListResponse
	userId := c.GetInt64("user_id")
	loggerFields := co.Service.Logger.GetLoggerFields(userId, c.ClientIP(), c.Request.Method, nil,
		c.Request.RequestURI)

	filter := make(map[string]interface{})
	filter["user_id"] = userId
	filter["category"] = cockarticletype.Draft
	// 草稿沒有收藏功能，userID為0
	cocktails, total, err := co.CocktailUsecase.GetAllWithFilter(c, filter, domain.PaginationUsecase{}, 0)
	if err != nil {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
			"get cocktails with filter failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	cocktailList := make([]viewmodels.DraftCocktailList, 0)
	for _, cocktail := range cocktails {
		out := viewmodels.DraftCocktailList{
			CocktailID:  cocktail.CocktailID,
			Title:       cocktail.Title,
			Photo:       cocktail.CoverPhoto.Photo,
			Description: cocktail.Description,
			CreatedDate: cocktail.CreatedDate,
		}

		cocktailList = append(cocktailList, out)
	}

	response.Total = total
	response.DraftCocktailList = cocktailList

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation POST /cocktails cocktail postArticleRequest
// ---
// summary: Post cocktail article.
// description: Post cocktail article.
//
// security:
// - Bearer: [apiKey]
//
// responses:
//   "201":
//     description: success

func (co *CocktailHandler) PostArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")
	userType := c.GetInt("type")

	var request viewmodels.PostArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		co.Service.Logger.LogFile(c, logrus.InfoLevel, co.Service.Logger.GetLoggerFields(userId, c.ClientIP(),
			c.Request.Method, nil, c.Request.RequestURI),
			"parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}
	loggerFields := co.Service.Logger.GetLoggerFields(userId, c.ClientIP(), c.Request.Method, request,
		c.Request.RequestURI)

	if userType == usertype.Test.Int() {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
			"!!test!! store article failed - %s", errors.New("test fail"))
		util.PackResponseWithError(c, domain.ErrInternalError, domain.ErrInternalError.Error())
		return
	}

	var cocktail = domain.Cocktail{
		UserID:      userId,
		Title:       request.Name,
		Description: request.Description,
		Category:    cockarticletype.Formal.Int(),
	}

	var ingredients []domain.CocktailIngredient
	for _, ingredient := range request.IngredientList {
		out := domain.CocktailIngredient{
			IngredientName:   ingredient.Name,
			IngredientAmount: ingredient.Amount,
		}
		ingredients = append(ingredients, out)
	}

	var steps []domain.CocktailStep
	for stepNumber, step := range request.StepList {
		out := domain.CocktailStep{
			StepNumber:      stepNumber,
			StepDescription: step.Description,
		}
		steps = append(steps, out)
	}

	var images []domain.CocktailImage
	for idx, file := range request.Files {
		var isCoverPhoto bool
		if idx == 0 {
			isCoverPhoto = true
		} else {
			isCoverPhoto = false
		}

		dataURL, err := dataurl.DecodeString(file)
		if err != nil {
			co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
				"decode data url failed - %s", err)
			util.PackResponseWithError(c, err, err.Error())
			return
		}

		out := domain.CocktailImage{
			File:         string(dataURL.Data),
			ContentType:  dataURL.MediaType.ContentType(),
			IsCoverPhoto: isCoverPhoto,
		}
		images = append(images, out)
	}

	err := co.CocktailUsecase.Store(c, &cocktail, ingredients, steps, images, userId)
	if err != nil {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
			"store article failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	util.PackResponseWithData(c, http.StatusCreated, nil, domain.GetErrorCode(nil), "")
}

// swagger:operation POST /cocktail-drafts cocktail postDraftArticleRequest
// ---
// summary: Post cocktail draft article.
// description: Post cocktail draft article.
//
// security:
// - Bearer: [apiKey]
//
// responses:
//   "201":
//     description: success

func (co *CocktailHandler) PostDraftArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")
	userType := c.GetInt("type")

	var request viewmodels.PostDraftArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		co.Service.Logger.LogFile(c, logrus.InfoLevel, co.Service.Logger.GetLoggerFields(userId, c.ClientIP(),
			c.Request.Method, nil, c.Request.RequestURI), "parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}
	loggerFields := co.Service.Logger.GetLoggerFields(userId, c.ClientIP(), c.Request.Method, request,
		c.Request.RequestURI)

	if userType == usertype.Test.Int() {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
			"!!test!! store article failed - %s", errors.New("test fail"))
		util.PackResponseWithError(c, domain.ErrInternalError, domain.ErrInternalError.Error())
		return
	}

	var cocktail = domain.Cocktail{
		UserID:      userId,
		Title:       request.Name,
		Description: request.Description,
		Category:    cockarticletype.Draft.Int(),
	}

	var ingredients []domain.CocktailIngredient
	for _, ingredient := range request.IngredientList {
		out := domain.CocktailIngredient{
			IngredientName:   ingredient.Name,
			IngredientAmount: ingredient.Amount,
		}
		ingredients = append(ingredients, out)
	}

	var steps []domain.CocktailStep
	for stepNumber, step := range request.StepList {
		out := domain.CocktailStep{
			StepNumber:      stepNumber,
			StepDescription: step.Description,
		}
		steps = append(steps, out)
	}

	var images []domain.CocktailImage
	for idx, file := range request.Files {
		var isCoverPhoto bool
		if idx == 0 {
			isCoverPhoto = true
		} else {
			isCoverPhoto = false
		}

		dataURL, err := dataurl.DecodeString(file)
		if err != nil {
			co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
				"decode data url failed - %s", err)
			util.PackResponseWithError(c, err, err.Error())
			return
		}

		out := domain.CocktailImage{
			File:         string(dataURL.Data),
			ContentType:  dataURL.MediaType.ContentType(),
			IsCoverPhoto: isCoverPhoto,
		}
		images = append(images, out)
	}

	err := co.CocktailUsecase.Store(c, &cocktail, ingredients, steps, images, userId)
	if err != nil {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
			"store article failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	util.PackResponseWithData(c, http.StatusCreated, nil, domain.GetErrorCode(nil), "")
}

// swagger:operation POST /cocktail-drafts/{id} cocktail makeCocktailDraftToFormal
// ---
// summary: Make cocktail draft article to formal article.
// description: Make cocktail draft article to formal article.
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
//   "201":
//     description: success

func (co *CocktailHandler) MakeDraftArticleToFormalArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")
	var request viewmodels.MakeDraftArticleToFormalArticle
	if err := c.ShouldBindUri(&request); err != nil {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, co.Service.Logger.GetLoggerFields(userId, c.ClientIP(),
			c.Request.Method, nil, c.Request.RequestURI), "parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}
	loggerFields := co.Service.Logger.GetLoggerFields(userId, c.ClientIP(), c.Request.Method, request,
		c.Request.RequestURI)

	err := co.CocktailUsecase.MakeDraftToFormal(c, request.ID, userId)
	if err != nil {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
			"make draft to formal failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	util.PackResponseWithData(c, http.StatusCreated, nil, domain.GetErrorCode(nil), "")
}

// swagger:operation PUT /cocktail-drafts/{id} cocktail updateCocktailDraft
// ---
// summary: Edit cocktail draft article.
// description: Edit cocktail draft article.
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
// - name: Body
//   in: body
//   schema:
//     "$ref": "#/definitions/updateDraftArticleRequest"
//
// responses:
//   "201":
//     description: success

func (co *CocktailHandler) UpdateDraftArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")

	var requestUri viewmodels.UpdateDraftArticleUriRequest
	if err := c.ShouldBindUri(&requestUri); err != nil {
		co.Service.Logger.LogFile(c, logrus.InfoLevel, co.Service.Logger.GetLoggerFields(userId, c.ClientIP(),
			c.Request.Method, nil, c.Request.RequestURI), "parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	var requestBody viewmodels.UpdateDraftArticleRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		co.Service.Logger.LogFile(c, logrus.InfoLevel, co.Service.Logger.GetLoggerFields(userId, c.ClientIP(),
			c.Request.Method, nil, c.Request.RequestURI), "parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	loggerFields := co.Service.Logger.GetLoggerFields(userId, c.ClientIP(), c.Request.Method, requestBody,
		c.Request.RequestURI)

	userType := c.GetInt("type")
	if userType == usertype.Test.Int() {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
			"!!test!! store article failed - %s", errors.New("test fail"))
		util.PackResponseWithError(c, domain.ErrInternalError, domain.ErrInternalError.Error())
		return
	}

	var cocktail = domain.Cocktail{
		CocktailID:  requestUri.ID,
		Title:       requestBody.Name,
		Description: requestBody.Description,
		Category:    cockarticletype.Draft.Int(),
	}

	var ingredients []domain.CocktailIngredient
	for _, ingredient := range requestBody.IngredientList {
		out := domain.CocktailIngredient{
			CocktailID:       requestUri.ID,
			IngredientName:   ingredient.Name,
			IngredientAmount: ingredient.Amount,
		}
		ingredients = append(ingredients, out)
	}

	var steps []domain.CocktailStep
	for stepNumber, step := range requestBody.StepList {
		out := domain.CocktailStep{
			CocktailID:      requestUri.ID,
			StepNumber:      stepNumber,
			StepDescription: step.Description,
		}
		steps = append(steps, out)
	}

	var images []domain.CocktailImage
	for idx, photo := range requestBody.Photos {
		var isCoverPhoto bool
		if idx == 0 {
			isCoverPhoto = true
		} else {
			isCoverPhoto = false
		}

		if photo.ImageFile != "" {
			dataURL, err := dataurl.DecodeString(photo.ImageFile)
			if err != nil {
				co.Service.Logger.LogFile(c, logrus.InfoLevel, loggerFields,
					"query by cocktail id failed - %s", err)
				util.PackResponseWithError(c, err, err.Error())
				return
			}

			out := domain.CocktailImage{
				ImageID:      photo.ID,
				CocktailID:   requestUri.ID,
				File:         string(dataURL.Data),
				ContentType:  dataURL.MediaType.ContentType(),
				IsCoverPhoto: isCoverPhoto,
				Order:        idx,
			}
			images = append(images, out)
		} else {
			out := domain.CocktailImage{
				ImageID:      photo.ID,
				CocktailID:   requestUri.ID,
				IsCoverPhoto: isCoverPhoto,
				Order:        idx,
			}
			images = append(images, out)
		}
	}

	err := co.CocktailUsecase.Update(c, &cocktail, ingredients, steps, images, userId)
	if err != nil {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields, "query by cocktail id failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	util.PackResponseWithData(c, http.StatusCreated, nil, domain.GetErrorCode(nil), "")
}

// swagger:operation PUT /cocktails/{id} cocktail updateFormalArticle
// ---
// summary: Edit cocktail formal article.
// description: Edit cocktail formal article.
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
// - name: Body
//   in: body
//   schema:
//     "$ref": "#/definitions/updateFormalArticleRequest"
//
// responses:
//
//	"200":
//	  description: success

func (co *CocktailHandler) UpdateFormalArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")

	var requestUri viewmodels.UpdateFormalArticleUriRequest
	if err := c.ShouldBindUri(&requestUri); err != nil {
		co.Service.Logger.LogFile(c, logrus.InfoLevel, co.Service.Logger.GetLoggerFields(userId, c.ClientIP(),
			c.Request.Method, nil, c.Request.RequestURI), "parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	var requestBody viewmodels.UpdateFormalArticleRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		co.Service.Logger.LogFile(c, logrus.InfoLevel, co.Service.Logger.GetLoggerFields(userId, c.ClientIP(),
			c.Request.Method, nil, c.Request.RequestURI), "parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	loggerFields := co.Service.Logger.GetLoggerFields(userId, c.ClientIP(), c.Request.Method, requestBody,
		c.Request.RequestURI)

	userType := c.GetInt("type")
	if userType == usertype.Test.Int() {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
			"!!test!! store article failed - %s", errors.New("test fail"))
		util.PackResponseWithError(c, domain.ErrInternalError, domain.ErrInternalError.Error())
		return
	}

	var cocktail = domain.Cocktail{
		CocktailID:  requestUri.ID,
		Title:       requestBody.Name,
		Description: requestBody.Description,
		Category:    cockarticletype.Formal.Int(),
	}

	var ingredients []domain.CocktailIngredient
	for _, ingredient := range requestBody.IngredientList {
		out := domain.CocktailIngredient{
			CocktailID:       requestUri.ID,
			IngredientName:   ingredient.Name,
			IngredientAmount: ingredient.Amount,
		}
		ingredients = append(ingredients, out)
	}

	var steps []domain.CocktailStep
	for stepNumber, step := range requestBody.StepList {
		out := domain.CocktailStep{
			CocktailID:      requestUri.ID,
			StepNumber:      stepNumber,
			StepDescription: step.Description,
		}
		steps = append(steps, out)
	}

	var images []domain.CocktailImage
	for idx, photo := range requestBody.Photos {
		var isCoverPhoto bool
		if idx == 0 {
			isCoverPhoto = true
		} else {
			isCoverPhoto = false
		}

		if photo.ImageFile != "" {
			dataURL, err := dataurl.DecodeString(photo.ImageFile)
			if err != nil {
				co.Service.Logger.LogFile(c, logrus.InfoLevel, loggerFields,
					"decode dataurl failed - %s", err)
				util.PackResponseWithError(c, err, err.Error())
				return
			}

			out := domain.CocktailImage{
				ImageID:      photo.ID,
				CocktailID:   requestUri.ID,
				File:         string(dataURL.Data),
				ContentType:  dataURL.MediaType.ContentType(),
				IsCoverPhoto: isCoverPhoto,
				Order:        idx,
			}
			images = append(images, out)
		} else {
			out := domain.CocktailImage{
				ImageID:      photo.ID,
				CocktailID:   requestUri.ID,
				IsCoverPhoto: isCoverPhoto,
				Order:        idx,
			}
			images = append(images, out)
		}
	}

	err := co.CocktailUsecase.Update(c, &cocktail, ingredients, steps, images, userId)
	if err != nil {
		co.Service.Logger.LogFile(c, logrus.InfoLevel, loggerFields, "update cocktail failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	util.PackResponseWithData(c, http.StatusCreated, nil, domain.GetErrorCode(nil), "")
}

// swagger:operation PATCH /cocktail-drafts cocktail deleteDraftArticleRequest
// ---
// summary: DELETE cocktail draft article.
// description: DELETE cocktail draft article.
//
// security:
// - Bearer: [apiKey]
//
// responses:
//   "200":
//     description: success

func (co *CocktailHandler) DeleteDraftArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")

	var request viewmodels.DeleteDraftArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		co.Service.Logger.LogFile(c, logrus.InfoLevel, co.Service.Logger.GetLoggerFields(userId, c.ClientIP(),
			c.Request.Method, nil, c.Request.RequestURI), "query by cocktail id failed - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	loggerFields := co.Service.Logger.GetLoggerFields(userId, c.ClientIP(), c.Request.Method, request,
		c.Request.RequestURI)

	userType := c.GetInt("type")
	if userType == usertype.Test.Int() {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
			"!!test!! store article failed - %s", errors.New("test fail"))
		util.PackResponseWithError(c, domain.ErrInternalError, domain.ErrInternalError.Error())
		return
	}

	for _, ids := range request.DeletedIds {
		err := co.CocktailUsecase.Delete(c, ids, userId)
		if err != nil {
			co.Service.Logger.LogFile(c, logrus.InfoLevel, loggerFields, "query by cocktail id failed - %s", err)
			util.PackResponseWithError(c, err, err.Error())
			return
		}
	}

	util.PackResponseWithData(c, http.StatusCreated, nil, domain.GetErrorCode(nil), "")
}

// swagger:operation PATCH /cocktails cocktail deleteFormalArticleRequest
// ---
// summary: DELETE cocktail formal article.
// description: DELETE cocktail formal article.
//
// security:
// - Bearer: [apiKey]
//
// responses:
//   "201":
//     description: success

func (co *CocktailHandler) DeleteFormalArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")

	var request viewmodels.DeleteFormalArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		co.Service.Logger.LogFile(c, logrus.InfoLevel, co.Service.Logger.GetLoggerFields(userId, c.ClientIP(),
			c.Request.Method, nil, c.Request.RequestURI), "query by cocktail id failed - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	loggerFields := co.Service.Logger.GetLoggerFields(userId, c.ClientIP(), c.Request.Method, request,
		c.Request.RequestURI)

	userType := c.GetInt("type")
	if userType == usertype.Test.Int() {
		co.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields,
			"!!test!! store article failed - %s", errors.New("test fail"))
		util.PackResponseWithError(c, domain.ErrInternalError, domain.ErrInternalError.Error())
		return
	}

	for _, ids := range request.DeletedIds {
		err := co.CocktailUsecase.Delete(c, ids, userId)
		if err != nil {
			co.Service.Logger.LogFile(c, logrus.InfoLevel, loggerFields, "query by cocktail id failed - %s", err)
			util.PackResponseWithError(c, err, err.Error())
			return
		}
	}

	util.PackResponseWithData(c, http.StatusCreated, nil, domain.GetErrorCode(nil), "")
}
