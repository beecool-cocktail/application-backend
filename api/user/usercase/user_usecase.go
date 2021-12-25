package usercase

import (
	"context"
	"errors"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/beecool-cocktail/application-backend/util"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

func (u *userUsecase) Logout(ctx context.Context, id int64) (err error) {

	token := util.GenString(64)
	redisToken := domain.UserCache{
		Id: id,
		AccessToken: token,
	}
	if err := u.userRedisRepo.UpdateToken(ctx, &redisToken); err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (u *userUsecase) QueryById(ctx context.Context, id int64) (*domain.User, error) {

	user, err := u.userMySQLRepo.QueryById(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Error(err)
		return nil, domain.ErrUserNotFound
	} else if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return user, nil
}