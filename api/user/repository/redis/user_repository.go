package redis

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/fatih/structs"
	"github.com/go-redis/redis/v9"
	"strconv"
)

type userRedisRepository struct {
	redis *redis.Client
}

type userToken struct {
	AccessToken  string `structs:"access_token"`
	RefreshToken string `structs:"refresh_token"`
}

type userInfo struct {
	Name string `structs:"name"`
}

func NewRedisUserRepository(redis *redis.Client) domain.UserRedisRepository {
	return &userRedisRepository{redis}
}

func (u *userRedisRepository) Store(ctx context.Context, r *domain.UserCache) error {
	key := "user:user_id:" + strconv.FormatInt(r.Id, 10)
	u.redis.HMSet(ctx, key, structs.Map(r))

	return nil
}

func (u *userRedisRepository) UpdateToken(ctx context.Context, r *domain.UserCache) error {
	key := "user:user_id:" + strconv.FormatInt(r.Id, 10)

	token := userToken{
		AccessToken:  r.AccessToken,
		RefreshToken: r.RefreshToken,
	}

	u.redis.HMSet(ctx, key, structs.Map(token))

	return nil
}

func (u *userRedisRepository) UpdateBasicInfo(ctx context.Context, r *domain.UserCache) error {
	key := "user:user_id:" + strconv.FormatInt(r.Id, 10)

	token := userInfo{
		Name: r.Name,
	}

	u.redis.HMSet(ctx, key, structs.Map(token))

	return nil
}

func (u *userRedisRepository) QueryUserNameByID(ctx context.Context, id int64) (string, error) {
	key := "user:user_id:" + strconv.FormatInt(id, 10)

	value, err := u.redis.HGet(ctx, key, "name").Result()
	if err == redis.Nil {
		return value, err
	} else if err != nil {
		return value, err
	}
	return value, nil
}
