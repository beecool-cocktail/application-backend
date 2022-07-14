package command

import (
	"context"
	"errors"
	"github.com/beecool-cocktail/application-backend/domain"
)

const (
	FavoriteCocktailDelete = "favorite_cocktail_delete"
)

type Operator interface {
	Undo(ctx context.Context, command *domain.Command) error
}

type OperatorHandler struct {
	Operator map[string]Operator
}

func NewOperatorHandler() OperatorHandler {

	operator := make(map[string]Operator)

	return OperatorHandler{
		Operator: operator,
	}
}

func (o OperatorHandler) SetOperator(operateName string, operator Operator) {
	o.Operator[operateName] = operator
}

func (o OperatorHandler) GetOperator(operateName string) (Operator, error) {

	if _, ok := o.Operator[operateName]; !ok {
		return nil, errors.New("no operator")
	}

	return o.Operator[operateName], nil
}

type Command interface {
	Execute(ctx context.Context, command *domain.Command) error
}

type UndoCommand struct {
	Operator Operator
}

func (u *UndoCommand) Execute(ctx context.Context, command *domain.Command) error {
	return u.Operator.Undo(ctx, command)
}

type Request struct {
	Command Command
}

func (r Request) Send(ctx context.Context, command *domain.Command) error {
	return r.Command.Execute(ctx, command)
}
