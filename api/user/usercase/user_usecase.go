package usercase

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/sirupsen/logrus"
)

type userUsecase struct {
	userMySQLRepo domain.UserMySQLRepository
	userRedisRepo domain.UserRedisRepository
}

func NewUserUsecase(clientMySQLRepo domain.UserMySQLRepository, clientRedisRepo domain.UserRedisRepository) domain.UserUsecase {
	return &userUsecase{
		userMySQLRepo: clientMySQLRepo,
		userRedisRepo: clientRedisRepo,
	}
}

func (c *userUsecase) Logout(ctx context.Context, id int64) (err error) {

	token := util.GenString(64)
	redisToken := domain.UserCache{
		Id: id,
		AccessToken: token,
	}
	if err := c.userRedisRepo.UpdateToken(ctx, &redisToken); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}