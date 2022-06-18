package http

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/enum/cockarticletype"
	"github.com/beecool-cocktail/application-backend/middleware"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/gin-gonic/gin"
	"github.com/vincent-petithory/dataurl"
	"net/http"
	"strconv"
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

	s.HTTP.GET("/api/cocktails/:cocktailID", middlewareHandler.JWTAuthMiddlewareIfExist(), handler.GetCocktailByCocktailID)
	s.HTTP.GET("/api/cocktail-drafts/:cocktailID", middlewareHandler.JWTAuthMiddleware(), handler.GetCocktailDraftByCocktailID)
	s.HTTP.GET("/api/cocktails", middlewareHandler.JWTAuthMiddlewareIfExist(), handler.CocktailList)
	s.HTTP.GET("/api/cocktail-drafts", middlewareHandler.JWTAuthMiddleware(), handler.CocktailDraftList)
	s.HTTP.POST("/api/cocktails", middlewareHandler.JWTAuthMiddleware(), handler.PostArticle)
	s.HTTP.POST("/api/cocktail-drafts", middlewareHandler.JWTAuthMiddleware(), handler.PostDraftArticle)
	s.HTTP.POST("/api/cocktail-drafts/:cocktailID", middlewareHandler.JWTAuthMiddleware(), handler.MakeDraftArticleToFormalArticle)
	s.HTTP.PUT("/api/cocktails/:cocktailID", middlewareHandler.JWTAuthMiddleware(), handler.UpdateFormalArticle)
	s.HTTP.PUT("/api/cocktail-drafts/:cocktailID", middlewareHandler.JWTAuthMiddleware(), handler.UpdateDraftArticle)
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
	var response viewmodels.GetCocktailByIDResponse
	userId := c.GetInt64("user_id")
	cocktailID := c.Param("cocktailID")
	api := "/cocktails" + cocktailID
	cocktailIDNumber, err := strconv.ParseInt(cocktailID, 10, 64)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	cocktail, err := co.CocktailUsecase.QueryByCocktailID(c, cocktailIDNumber, userId)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("query by cocktail id failed - %s", err)
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

	photos := make([]viewmodels.CocktailPhotoWithIDInResponse, 0)
	for _, photo := range cocktail.Photos {
		out := viewmodels.CocktailPhotoWithIDInResponse{
			ID:    photo.ID,
			Photo: photo.Photo,
		}
		photos = append(photos, out)
	}

	cocktailUser, err := co.UserUsecase.QueryById(c, cocktail.UserID)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("query user by user id failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	response = viewmodels.GetCocktailByIDResponse{
		CocktailID:     cocktail.CocktailID,
		UserID:         cocktail.UserID,
		UserName:       cocktail.UserName,
		UserPhoto:      cocktailUser.Photo,
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
	var response viewmodels.GetCocktailDraftByIDResponse
	cocktailID := c.Param("cocktailID")
	userId := c.GetInt64("user_id")

	api := "/cocktail-drafts" + cocktailID
	cocktailIDNumber, err := strconv.ParseInt(cocktailID, 10, 64)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	cocktail, err := co.CocktailUsecase.QueryDraftByCocktailID(c, cocktailIDNumber, userId)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("query by cocktail id failed - %s", err)
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

	photos := make([]viewmodels.CocktailPhotoWithIDInResponse, 0)
	for _, photo := range cocktail.Photos {
		out := viewmodels.CocktailPhotoWithIDInResponse{
			ID:    photo.ID,
			Photo: photo.Photo,
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
//    "$ref": "#/responses/popularCocktailListResponse"
func (co *CocktailHandler) CocktailList(c *gin.Context) {
	api := "/cocktails"
	userId := c.GetInt64("user_id")

	var response viewmodels.GetPopularCocktailListResponse
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	keyword := c.DefaultQuery("keyword", "")

	var cocktails []domain.APICocktail
	var total int64
	if keyword != "" && co.Service.Configure.Elastic.Enable {
		cocktails, total, err = co.CocktailUsecase.Search(c, keyword, page, pageSize, userId)
		if err != nil {
			service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("get cocktails with keyword "+
				"failed - %s", err)
			util.PackResponseWithError(c, err, err.Error())
			return
		}
	} else {
		filter := make(map[string]interface{})
		filter["category"] = cockarticletype.Formal
		cocktails, total, err = co.CocktailUsecase.GetAllWithFilter(c, filter, domain.PaginationUsecase{
			Page:     page,
			PageSize: pageSize,
		},
			userId)
		if err != nil {
			service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("get cocktails with filter failed - %s", err)
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

		photos := make([]string, 0)
		for _, photo := range cocktail.Photos {
			photos = append(photos, photo.Photo)
		}

		lowQualityPhotos := make([]string, 0)
		for _, photo := range cocktail.LowQualityPhotos {
			lowQualityPhotos = append(lowQualityPhotos, photo.Photo)
		}

		out := viewmodels.PopularCocktailList{
			CocktailID:       cocktail.CocktailID,
			UserID:           cocktail.UserID,
			UserName:         cocktail.UserName,
			Title:            cocktail.Title,
			Photos:           photos,
			LowQualityPhotos: lowQualityPhotos,
			IngredientList:   ingredients,
			IsCollected:      cocktail.IsCollected,
			CreatedDate:      cocktail.CreatedDate,
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
//  "200":
//    "$ref": "#/responses/getDraftCocktailListResponse"
func (co *CocktailHandler) CocktailDraftList(c *gin.Context) {
	api := "/cocktails"
	var response viewmodels.GetDraftCocktailListResponse
	userId := c.GetInt64("user_id")

	filter := make(map[string]interface{})
	filter["user_id"] = userId
	filter["category"] = cockarticletype.Draft
	// 草稿沒有收藏功能，userID為0
	cocktails, total, err := co.CocktailUsecase.GetAllWithFilter(c, filter, domain.PaginationUsecase{}, 0)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("get cocktails with filter failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	cocktailList := make([]viewmodels.DraftCocktailList, 0)
	for _, cocktail := range cocktails {
		out := viewmodels.DraftCocktailList{
			CocktailID:  cocktail.CocktailID,
			Title:       cocktail.Title,
			Photo:       cocktail.CoverPhoto.Photo,
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
//  "201": success
func (co *CocktailHandler) PostArticle(c *gin.Context) {
	api := "cocktail"
	userId := c.GetInt64("user_id")

	var request viewmodels.PostArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, request).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
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
			service.GetLoggerEntry(co.Service.Logger, api, request).Errorf("decode data url failed - %s", err)
			util.PackResponseWithError(c, err, err.Error())
			return
		}

		out := domain.CocktailImage{
			Data:         string(dataURL.Data),
			Type:         dataURL.MediaType.ContentType(),
			IsCoverPhoto: isCoverPhoto,
		}
		images = append(images, out)
	}

	err := co.CocktailUsecase.Store(c, &cocktail, ingredients, steps, images, userId)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, request).Errorf("store article failed - %s", err)
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
//  "201": success
func (co *CocktailHandler) PostDraftArticle(c *gin.Context) {
	api := "cocktail"
	userId := c.GetInt64("user_id")

	var request viewmodels.PostDraftArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, request).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
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
			service.GetLoggerEntry(co.Service.Logger, api, request).Errorf("decode data url failed - %s", err)
			util.PackResponseWithError(c, err, err.Error())
			return
		}

		out := domain.CocktailImage{
			Data:         string(dataURL.Data),
			Type:         dataURL.MediaType.ContentType(),
			IsCoverPhoto: isCoverPhoto,
		}
		images = append(images, out)
	}

	err := co.CocktailUsecase.Store(c, &cocktail, ingredients, steps, images, userId)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, request).Errorf("store article failed - %s", err)
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
//  "200": success
func (co *CocktailHandler) MakeDraftArticleToFormalArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")
	cocktailID := c.Param("cocktailID")
	api := "POST /cocktail-drafts/" + cocktailID

	cocktailIDNumber, err := strconv.ParseInt(cocktailID, 10, 64)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	err = co.CocktailUsecase.MakeDraftToFormal(c, cocktailIDNumber, userId)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("make draft to formal failed - %s", err)
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
//  "200": success
func (co *CocktailHandler) UpdateDraftArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")
	cocktailID := c.Param("cocktailID")
	api := "PUT /cocktail-drafts/" + cocktailID

	cocktailIDNumber, err := strconv.ParseInt(cocktailID, 10, 64)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	var request viewmodels.UpdateDraftArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, request).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	var cocktail = domain.Cocktail{
		CocktailID:  cocktailIDNumber,
		Title:       request.Name,
		Description: request.Description,
	}

	var ingredients []domain.CocktailIngredient
	for _, ingredient := range request.IngredientList {
		out := domain.CocktailIngredient{
			CocktailID:       cocktailIDNumber,
			IngredientName:   ingredient.Name,
			IngredientAmount: ingredient.Amount,
		}
		ingredients = append(ingredients, out)
	}

	var steps []domain.CocktailStep
	for stepNumber, step := range request.StepList {
		out := domain.CocktailStep{
			CocktailID:      cocktailIDNumber,
			StepNumber:      stepNumber,
			StepDescription: step.Description,
		}
		steps = append(steps, out)
	}

	var images []domain.CocktailImage
	for idx, photo := range request.Photos {
		var isCoverPhoto bool
		if idx == 0 {
			isCoverPhoto = true
		} else {
			isCoverPhoto = false
		}

		if photo.Photo != "" {
			dataURL, err := dataurl.DecodeString(photo.Photo)
			if err != nil {
				service.GetLoggerEntry(co.Service.Logger, api, request).Errorf("decode data url failed - %s", err)
				util.PackResponseWithError(c, err, err.Error())
				return
			}

			out := domain.CocktailImage{
				ImageID:      photo.ID,
				CocktailID:   cocktailIDNumber,
				Data:         string(dataURL.Data),
				Type:         dataURL.MediaType.ContentType(),
				IsCoverPhoto: isCoverPhoto,
			}
			images = append(images, out)
		} else {
			out := domain.CocktailImage{
				ImageID:      photo.ID,
				CocktailID:   cocktailIDNumber,
				IsCoverPhoto: isCoverPhoto,
			}
			images = append(images, out)
		}
	}

	err = co.CocktailUsecase.Update(c, &cocktail, ingredients, steps, images, userId)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("update draft cocktail failed - %s", err)
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
//  "200": success
func (co *CocktailHandler) UpdateFormalArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")
	cocktailID := c.Param("cocktailID")
	api := "PUT /cocktails/" + cocktailID

	cocktailIDNumber, err := strconv.ParseInt(cocktailID, 10, 64)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	var request viewmodels.UpdateFormalArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, request).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	var cocktail = domain.Cocktail{
		CocktailID:  cocktailIDNumber,
		Title:       request.Name,
		Description: request.Description,
	}

	var ingredients []domain.CocktailIngredient
	for _, ingredient := range request.IngredientList {
		out := domain.CocktailIngredient{
			CocktailID:       cocktailIDNumber,
			IngredientName:   ingredient.Name,
			IngredientAmount: ingredient.Amount,
		}
		ingredients = append(ingredients, out)
	}

	var steps []domain.CocktailStep
	for stepNumber, step := range request.StepList {
		out := domain.CocktailStep{
			CocktailID:      cocktailIDNumber,
			StepNumber:      stepNumber,
			StepDescription: step.Description,
		}
		steps = append(steps, out)
	}

	var images []domain.CocktailImage
	for idx, photo := range request.Photos {
		var isCoverPhoto bool
		if idx == 0 {
			isCoverPhoto = true
		} else {
			isCoverPhoto = false
		}

		if photo.Photo != "" {
			dataURL, err := dataurl.DecodeString(photo.Photo)
			if err != nil {
				service.GetLoggerEntry(co.Service.Logger, api, request).Errorf("decode data url failed - %s", err)
				util.PackResponseWithError(c, err, err.Error())
				return
			}

			out := domain.CocktailImage{
				ImageID:      photo.ID,
				CocktailID:   cocktailIDNumber,
				Data:         string(dataURL.Data),
				Type:         dataURL.MediaType.ContentType(),
				IsCoverPhoto: isCoverPhoto,
			}
			images = append(images, out)
		} else {
			out := domain.CocktailImage{
				ImageID:      photo.ID,
				CocktailID:   cocktailIDNumber,
				IsCoverPhoto: isCoverPhoto,
			}
			images = append(images, out)
		}
	}

	err = co.CocktailUsecase.Update(c, &cocktail, ingredients, steps, images, userId)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("update draft cocktail failed - %s", err)
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
//  "200": success
func (co *CocktailHandler) DeleteDraftArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")
	cocktailID := c.Param("cocktailID")
	api := "DELETE /cocktail-drafts/" + cocktailID

	var request viewmodels.DeleteDraftArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, request).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	for _, ids := range request.DeletedIds {
		err := co.CocktailUsecase.Delete(c, ids, userId)
		if err != nil {
			service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("delete draft cocktail failed - %s", err)
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
//  "200": success
func (co *CocktailHandler) DeleteFormalArticle(c *gin.Context) {
	userId := c.GetInt64("user_id")
	api := "DELETE /cocktails"

	var request viewmodels.DeleteFormalArticleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, request).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, domain.ErrParameterIllegal, domain.ErrParameterIllegal.Error())
		return
	}

	for _, ids := range request.DeletedIds {
		err := co.CocktailUsecase.Delete(c, ids, userId)
		if err != nil {
			service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("delete draft cocktail failed - %s", err)
			util.PackResponseWithError(c, err, err.Error())
			return
		}
	}

	util.PackResponseWithData(c, http.StatusCreated, nil, domain.GetErrorCode(nil), "")
}
