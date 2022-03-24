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
}

func NewCocktailHandler(s *service.Service, cocktailUsecase domain.CocktailUsecase, middlewareHandler middleware.Handler) {
	handler := &CocktailHandler{
		Service:         s,
		CocktailUsecase: cocktailUsecase,
	}

	s.HTTP.GET("/api/cocktails/:cocktailID", handler.GetCocktailByCocktailID)
	s.HTTP.GET("/api/cocktail-drafts/:cocktailID", middlewareHandler.JWTAuthMiddleware(), handler.GetCocktailDraftByCocktailID)
	s.HTTP.GET("/api/cocktails", handler.CocktailList)
	s.HTTP.GET("/api/cocktail-drafts", middlewareHandler.JWTAuthMiddleware(), handler.CocktailDraftList)
	s.HTTP.POST("/api/cocktails", middlewareHandler.JWTAuthMiddleware(), handler.PostArticle)
	s.HTTP.POST("/api/cocktail-drafts", middlewareHandler.JWTAuthMiddleware(), handler.PostDraftArticle)
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
	cocktailID := c.Param("cocktailID")
	api := "/cocktails" + cocktailID
	cocktailIDNumber, err := strconv.ParseInt(cocktailID, 10, 64)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	cocktail, err := co.CocktailUsecase.QueryByCocktailID(c, cocktailIDNumber)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("query by cocktail id failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	ingredients := make([]viewmodels.CocktailIngredient, 0)
	for _, ingredient := range cocktail.Ingredients {
		out := viewmodels.CocktailIngredient{
			Name:   ingredient.IngredientName,
			Amount: ingredient.IngredientAmount,
			Unit:   ingredient.IngredientUnit,
		}
		ingredients = append(ingredients, out)
	}

	steps := make([]viewmodels.CocktailStep, 0)
	for _, step := range cocktail.Steps {
		out := viewmodels.CocktailStep{
			Description: step.StepDescription,
		}
		steps = append(steps, out)
	}

	response = viewmodels.GetCocktailByIDResponse{
		CocktailID:     cocktail.CocktailID,
		UserID:         cocktail.UserID,
		UserName:       cocktail.UserName,
		Photos:         cocktail.Photos,
		Title:          cocktail.Title,
		Description:    cocktail.Description,
		IngredientList: ingredients,
		StepList:       steps,
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

	ingredients := make([]viewmodels.CocktailIngredient, 0)
	for _, ingredient := range cocktail.Ingredients {
		out := viewmodels.CocktailIngredient{
			Name:   ingredient.IngredientName,
			Amount: ingredient.IngredientAmount,
			Unit:   ingredient.IngredientUnit,
		}
		ingredients = append(ingredients, out)
	}

	steps := make([]viewmodels.CocktailStep, 0)
	for _, step := range cocktail.Steps {
		out := viewmodels.CocktailStep{
			Description: step.StepDescription,
		}
		steps = append(steps, out)
	}

	photos := make([]string, 0)
	for _, path := range cocktail.Photos {
		photos = append(photos, path)
	}

	response = viewmodels.GetCocktailDraftByIDResponse{
		CocktailID:     cocktail.CocktailID,
		Photos:         photos,
		Title:          cocktail.Title,
		Description:    cocktail.Description,
		IngredientList: ingredients,
		StepList:       steps,
		CreatedDate:    cocktail.CreatedDate,
	}

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation GET /cocktails cocktail noRequest
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
// responses:
//  "200":
//    "$ref": "#/responses/popularCocktailListResponse"
func (co *CocktailHandler) CocktailList(c *gin.Context) {
	api := "/cocktails"
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

	filter := make(map[string]interface{})
	filter["category"] = cockarticletype.Normal
	cocktails, total, err := co.CocktailUsecase.GetAllWithFilter(c, filter, domain.PaginationUsecase{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("get cocktails with filter failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	cocktailList := make([]viewmodels.PopularCocktailList, 0)
	for _, cocktail := range cocktails {
		ingredients := make([]viewmodels.CocktailIngredient, 0)
		for _, ingredient := range cocktail.Ingredients {
			out := viewmodels.CocktailIngredient{
				Name:   ingredient.IngredientName,
				Amount: ingredient.IngredientAmount,
				Unit:   ingredient.IngredientUnit,
			}
			ingredients = append(ingredients, out)
		}

		out := viewmodels.PopularCocktailList{
			CocktailID:     cocktail.CocktailID,
			UserID:         cocktail.UserID,
			UserName:       cocktail.UserName,
			Title:          cocktail.Title,
			Photos:         cocktail.Photos,
			IngredientList: ingredients,
			CreatedDate:    cocktail.CreatedDate,
		}

		cocktailList = append(cocktailList, out)
	}

	response.Total = total
	response.PopularCocktailList = cocktailList

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}

// swagger:operation GET /cocktail-drafts cocktail noRequest
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
	cocktails, total, err := co.CocktailUsecase.GetAllWithFilter(c, filter, domain.PaginationUsecase{})
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
			Photo:       cocktail.CoverPhoto,
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
		Category:    cockarticletype.Normal.Int(),
	}

	var ingredients []domain.CocktailIngredient
	for _, ingredient := range request.IngredientList {
		out := domain.CocktailIngredient{
			IngredientName:   ingredient.Name,
			IngredientAmount: ingredient.Amount,
			IngredientUnit:   ingredient.Unit,
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

	err := co.CocktailUsecase.Store(c, &cocktail, ingredients, steps, images)
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
			IngredientUnit:   ingredient.Unit,
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

	err := co.CocktailUsecase.Store(c, &cocktail, ingredients, steps, images)
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, request).Errorf("store article failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	util.PackResponseWithData(c, http.StatusCreated, nil, domain.GetErrorCode(nil), "")
}
