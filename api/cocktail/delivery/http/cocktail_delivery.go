package http

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CocktailHandler struct {
	Service *service.Service
	CocktailUsecase domain.CocktailUsecase
}

func NewCocktailHandler(s *service.Service, cocktailUsecase domain.CocktailUsecase) {
	handler := &CocktailHandler{
		Service: s,
		CocktailUsecase: cocktailUsecase,
	}

	s.HTTP.GET("/api/cocktails", handler.CocktailList)
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
		Page: page,
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
			CocktailID: cocktail.CocktailID,
			Photo: cocktail.Photo,
			Title: cocktail.Title,
			CreatedDate: util.GetFormatTime(cocktail.CreatedDate, "UTC"),
		}

		cocktailList = append(cocktailList, out)
	}

	response.Total = total
	response.PopularCocktailList = cocktailList

	util.PackResponseWithData(c, http.StatusOK, response, domain.GetErrorCode(nil), "")
}