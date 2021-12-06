package util

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/gin-gonic/gin"
)

func PackResponseWithData(c *gin.Context, httpStatusCode int, data interface{}, errorCode, errorMessage string) {
	response := viewmodels.ResponseData{
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
		Data:         data,
	}

	c.JSON(httpStatusCode, response)
}

func PackResponseWithError(c * gin.Context, error error, errorMessage string) {
	response := viewmodels.ResponseData{
		ErrorCode:    domain.GetErrorCode(error),
		ErrorMessage: errorMessage,
	}

	c.JSON(domain.GetStatusCode(error), response)
}
