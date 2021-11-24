package usercase

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
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

func (c *UserUsecase) Register(ctx context.Context, d *domain.User, r *domain.UserCache) (id int64, token string, err error) {

	token = util.GenString(64)
	id = util.GetID(util.IdGenerator)
	d.ID = id

	if err := c.userMySQLRepo.Store(ctx, d); err != nil {
		logrus.Error(err)
		return id, token, err
	}

	r.AccessToken = token
	r.TokenExpire = util.GetFormatTime(time.Now().Add(time.Duration(60) * time.Minute), "Local")
	key := "client:client_id:" + strconv.FormatInt(d.ID, 10)
	if err := c.userRedisRepo.Store(ctx, r, key); err != nil {
		logrus.Error(err)
		return id, token, err
	}

	return id, token, nil
}