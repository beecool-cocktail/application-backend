package util

import (
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/gin-gonic/gin"
)

func PackResponseWithData(c * gin.Context, httpStatusCode int, data interface{}, errorCode, errorMessage string) {
	response := viewmodels.ResponseData{
		ErrorCode: errorCode,
		ErrorMessage: errorMessage,
		Data: data,
	}

	c.JSON(httpStatusCode, response)
}

func PackResponseWithError(c * gin.Context, httpStatusCode int, errorCode, errorMessage string) {
	response := viewmodels.ResponseData{
		ErrorCode: errorCode,
		ErrorMessage: errorMessage,
	}

	c.JSON(httpStatusCode, response)
}