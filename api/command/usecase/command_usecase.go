package usecase

import (
	"context"
	cmdPattern "github.com/beecool-cocktail/application-backend/command"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/go-redis/redis"
)

type commandUsecase struct {
	service          *service.Service
	commandRedisRepo domain.CommandRedisRepository
	operatorHandler  cmdPattern.OperatorHandler
}

func NewCommandUsecase(s *service.Service, commandRedisRepo domain.CommandRedisRepository,
	operatorHandler cmdPattern.OperatorHandler) domain.CommandUsecase {
	return &commandUsecase{
		service:          s,
		commandRedisRepo: commandRedisRepo,
		operatorHandler:  operatorHandler,
	}
}

func (c *commandUsecase) Store(ctx context.Context, dc *domain.Command) error {

	err := c.commandRedisRepo.Store(ctx, dc)
	if err != nil {
		return err
	}

	return nil
}

func (c *commandUsecase) Undo(ctx context.Context, id string) error {

	command, err := c.commandRedisRepo.GetByID(ctx, id)
	if err == redis.Nil {
		return domain.ErrCommandNotFound
	} else if err != nil {
		return err
	}

	operator, err := c.operatorHandler.GetOperator(command.Name)
	if err != nil {
		return err
	}

	undoCommand := &cmdPattern.UndoCommand{
		Operator: operator,
	}

	undoRequest := &cmdPattern.Request{
		Command: undoCommand,
	}

	err = undoRequest.Send(ctx, &command)
	if err != nil {
		return err
	}

	err = c.commandRedisRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
