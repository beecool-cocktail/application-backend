package http

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CocktailHandler struct {
	CocktailUsecase domain.CocktailUsecase
}

func NewCocktailHandler(s *service.Service, cocktailUsecase domain.CocktailUsecase) {
	handler := &CocktailHandler{
		CocktailUsecase: cocktailUsecase,
	}

	s.HTTP.POST("/api/cocktails", handler.CocktailList)
}

// swagger:route POST /cocktails cocktail popularCocktailListRequest
//
// Get popular cocktail list
//
// Get popular cocktail list order by create date.
//
// security:
// - Bearer: []
//
// Responses:
//   200: popularCocktailListResponse
//   400: description: bad request
//   401: description: unauthorized
//   404: description: item not found
//   500: description: internal error
func (co *CocktailHandler) CocktailList(c *gin.Context) {

	var request viewmodels.GetPopularCocktailListRequest
	var response viewmodels.GetPopularCocktailListResponse
	if err := c.BindJSON(&request); err != nil {
		logrus.Error(err)
		util.PackResponseWithError(c, domain.ErrRequestDecodeFailed, "request unmarshal failed")
	}

	cocktails, total, err := co.CocktailUsecase.GetAllWithFilter(c, nil, domain.PaginationUsecase{
		Page: request.Page,
		PageSize: request.PageSize,
	})
	if err != nil {
		logrus.Error(err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	var cocktailList []viewmodels.PopularCocktailList
	for _, cocktail := range *cocktails {
		out := viewmodels.PopularCocktailList{
			CocktailID: cocktail.CocktailID,
			Photo: "",
			Title: cocktail.Title,
			CreatedDate: util.GetFormatTime(cocktail.CreatedDate, "UTC"),
		}

		cocktailList = append(cocktailList, out)
	}

	response.Total = total
	response.PopularCocktailList = cocktailList

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}