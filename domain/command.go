package domain

import (
	"context"
	"time"
)

type CreateCommand struct {
	OperatorID interface{} `json:"operator_id"`
	TargetID   interface{} `json:"target_id"`
}

type EditCommand struct {
	OperatorID interface{} `json:"operator_id"`
	TargetID   interface{} `json:"target_id"`
	Fields     []EditField `json:"fields"`
}

type EditField struct {
	Field  string      `json:"field"`
	Before interface{} `json:"before"`
	After  interface{} `json:"after"`
}

type DeleteCommand struct {
	OperatorID interface{} `json:"operator_id"`
	TargetID   interface{} `json:"target_id"`
}

type CommandType struct {
	Create CreateCommand `json:"create"`
	Edit   EditCommand   `json:"edit"`
	Delete DeleteCommand `json:"delete"`
}

type Command struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Type       CommandType   `json:"type"`
	ExpireTime time.Duration `json:"expire_time"`
}

type CommandRedisRepository interface {
	Store(ctx context.Context, c *Command) error
	GetByID(ctx context.Context, id string) (Command, error)
	Delete(ctx context.Context, id string) error
}

type CommandUsecase interface {
	Store(ctx context.Context, c *Command) error
	Undo(ctx context.Context, id string) error
}
