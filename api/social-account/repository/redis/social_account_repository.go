package redis

import (
	"context"
	"time"

	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/fatih/structs"
	"github.com/go-redis/redis/v9"
)

type socialAccountRedisRepository struct {
	redis *redis.Client
}

func NewRedisSocialAccountRepository(redis *redis.Client) domain.SocialAccountRedisRepository {
	return &socialAccountRedisRepository{redis}
}

func (s *socialAccountRedisRepository) StoreState(ctx context.Context, key string, state domain.State) error {

	s.redis.HMSet(ctx, key, structs.Map(state))
	s.redis.Expire(ctx, key, 5*60*time.Second)

	return nil
}

func (s *socialAccountRedisRepository) GetState(ctx context.Context, key string) (domain.State, error) {

	var state domain.State
	err := s.redis.HMGet(ctx, key, "redirect_path", "collect_after_login").Scan(&state)

	return state, err
}

func (s *socialAccountRedisRepository) DeleteState(ctx context.Context, key string) error {

	s.redis.Del(ctx, key)

	return nil
}
