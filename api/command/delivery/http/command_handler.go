package http

import (
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/middleware"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type CommandHandler struct {
	Configure      *service.Configure
	Logger         *logrus.Logger
	CommandUsecase domain.CommandUsecase
}

func NewCommandHandler(s *service.Service,
	commandUsecase domain.CommandUsecase,
	middlewareHandler middleware.Handler) {

	handler := &CommandHandler{
		Configure:      s.Configure,
		Logger:         s.Logger,
		CommandUsecase: commandUsecase,
	}

	s.HTTP.POST("/api/command/:commandID/undo", middlewareHandler.JWTAuthMiddleware(), handler.UndoCommand)
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
//  "200": success
func (ch *CommandHandler) UndoCommand(c *gin.Context) {
	commandID := c.Param("commandID")

	api := "/command/" + commandID + "/undo"

	err := ch.CommandUsecase.Undo(c, commandID)
	if err != nil {
		service.GetLoggerEntry(ch.Logger, api, nil).Errorf("undo command failed - %s", err)
		util.PackResponseWithError(c, err, err.Error())
		return
	}

	util.PackResponseWithData(c, http.StatusOK, nil, domain.GetErrorCode(nil), "")
}
