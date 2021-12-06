package usercase

import (
	"github.com/beecool-cocktail/application-backend/domain"
)

type UserUsecase struct {
	userMySQLRepo domain.UserMySQLRepository
	userRedisRepo domain.UserRedisRepository
}

func NewUserUsecase(clientMySQLRepo domain.UserMySQLRepository, clientRedisRepo domain.UserRedisRepository) domain.UserUsecase {
	return &UserUsecase{
		userMySQLRepo: clientMySQLRepo,
		userRedisRepo: clientRedisRepo,
	}
}
