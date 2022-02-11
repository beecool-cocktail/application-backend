package http

import (
	"fmt"
	"github.com/beecool-cocktail/application-backend/domain"
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

	s.HTTP.GET("/api/cocktails", handler.CocktailList)
	s.HTTP.POST("/api/cocktails", middlewareHandler.JWTAuthMiddleware(), handler.PostArticle)
}

// swagger:operation GET /cocktails cocktail noRequest
//
// summary: Get popular cocktail list
// description: Get popular cocktail list order by create date.
//
// ---
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

	cocktails, total, err := co.CocktailUsecase.GetAllWithFilter(c, nil, domain.PaginationUsecase{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		service.GetLoggerEntry(co.Service.Logger, api, nil).Errorf("get cocktails with filter failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	var cocktailList []viewmodels.PopularCocktailList
	for _, cocktail := range cocktails {
		out := viewmodels.PopularCocktailList{
			CocktailID:  cocktail.CocktailID,
			Title:       cocktail.Title,
			Photo:       cocktail.Photo,
			CreatedDate: cocktail.CreatedDate,
		}

		cocktailList = append(cocktailList, out)
	}

	response.Total = total
	response.PopularCocktailList = cocktailList

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
	fmt.Println(userId)

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
	for _, step := range request.StepList {
		out := domain.CocktailStep{
			StepNumber:      step.Step,
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
