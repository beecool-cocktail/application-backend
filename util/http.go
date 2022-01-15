package util

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/gin-gonic/gin"
	"reflect"
)

func PackResponseWithData(c *gin.Context, httpStatusCode int, data interface{}, errorCode, errorMessage string) {
	var res viewmodels.ResponseData
	emptySlice := make([]interface{}, 0)
	if val := reflect.ValueOf(data);
		val.Kind() == reflect.Slice && val.Len() == 0 || data == nil{
		res = viewmodels.ResponseData{
			ErrorCode:    errorCode,
			ErrorMessage: errorMessage,
			Data:         emptySlice,
		}
	} else {
		res = viewmodels.ResponseData{
			ErrorCode:    errorCode,
			ErrorMessage: errorMessage,
			Data:         data,
		}
	}

	c.JSON(httpStatusCode, res)
}

func PackResponseWithError(c *gin.Context, error error, errorMessage string) {
	emptySlice := make([]interface{}, 0)
	response := viewmodels.ResponseData{
		ErrorCode:    domain.GetErrorCode(error),
		ErrorMessage: errorMessage,
		Data:         emptySlice,
	}

	c.JSON(domain.GetStatusCode(error), response)
}