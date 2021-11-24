package redis

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/fatih/structs"
	"github.com/go-redis/redis"
)

type userRedisRepository struct {
	redis *redis.Client
}

func NewRedisUserRepository(redis *redis.Client) domain.UserRedisRepository {
	return &userRedisRepository{redis}
}

func (u *userRedisRepository) Store(ctx context.Context, r *domain.UserCache, key string) error {
	u.redis.HMSet(key, structs.Map(r))

	return nil
}