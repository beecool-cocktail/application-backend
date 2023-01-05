package http

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/middleware"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/beecool-cocktail/application-backend/viewmodels"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CommandHandler struct {
	Service        *service.Service
	CommandUsecase domain.CommandUsecase
}

func NewCommandHandler(s *service.Service,
	commandUsecase domain.CommandUsecase,
	middlewareHandler middleware.Handler) {

	handler := &CommandHandler{
		Service:        s,
		CommandUsecase: commandUsecase,
	}

	s.HTTP.POST("/api/command/:id/undo", middlewareHandler.JWTAuthMiddleware(), handler.UndoCommand)
}

// swagger:operation POST /command/{id}/undo command undoCommand
// ---
// summary: Undo command.
// description: Undo command.
//
// security:
// - Bearer: [apiKey]
//
// parameters:
// - name: id
//   in: path
//   required: true
//   type: string
//   example: 1
//
// responses:
//   "200":
//     description: success

func (ch *CommandHandler) UndoCommand(c *gin.Context) {
	userId := c.GetInt64("user_id")
	var request viewmodels.CommandUndoRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		ch.Service.Logger.LogFile(c, logrus.InfoLevel, ch.Service.Logger.GetLoggerFields(userId, c.ClientIP(),
			c.Request.Method, nil, c.Request.RequestURI), "parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
	}

	loggerFields := ch.Service.Logger.GetLoggerFields(userId, c.ClientIP(), c.Request.Method, nil,
		c.Request.RequestURI)

	err = ch.CommandUsecase.Undo(c, request.ID)
	if err != nil {
		ch.Service.Logger.LogFile(c, logrus.ErrorLevel, loggerFields, "parameter illegal - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	util.PackResponseWithData(c, http.StatusOK, nil, domain.GetErrorCode(nil), "")
}
