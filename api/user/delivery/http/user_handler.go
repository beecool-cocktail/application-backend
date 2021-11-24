package http

import (
	"encoding/json"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/middleware"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

func NewUserHandler(e *gin.Engine, userUsecase domain.UserUsecase, middlewareHandler *middleware.Handler) {
	handler := &UserHandler{
		UserUsecase: userUsecase,
	}

	e.POST("/api/client/register", handler.ClientRegister)
}

// swagger:route POST /client/register client clientRegisterRequest
//
// register.
//
// register a client.
//
//     Responses:
//       201: clientRegisterResponse
func (d *UserHandler) ClientRegister(c *gin.Context) {
	var request viewmodels.ClientRegisterRequest
	var response viewmodels.ClientRegisterResponse
	if err := c.BindJSON(&request); err != nil {
		logrus.Error(err)
		util.PackResponseWithError(c, http.StatusBadRequest, "", "request unmarshal failed")
	}

	id, token, err := d.UserUsecase.Register(c,
		&domain.User{
			Account:  request.Account,
			Password: request.Password,
		}, &domain.UserCache{
			Account: request.Account,
		})
	if err != nil {
		logrus.Error(err)
		util.PackResponseWithError(c, http.StatusBadRequest, "", "Internal error. Store failed")
		return
	}

	response = viewmodels.ClientRegisterResponse{
		Token: token,
	}

	result, err := json.Marshal(response)
	if err != nil {
		logrus.Error("json marshal failed")
		util.PackResponseWithError(c, http.StatusBadRequest, "", "json marshal error")
	}

	session := sessions.Default(c)
	session.Set("account", request.Account)
	session.Set("id", id)
	session.Save()

	util.PackResponseWithData(c, http.StatusCreated, result, "", "")
}