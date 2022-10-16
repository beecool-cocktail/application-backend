package redis

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v9"
	"time"
)

type cocktailRedisRepository struct {
	redis *redis.Client
}

func NewRedisCocktailRepository(redis *redis.Client) domain.CocktailRedisRepository {
	return &cocktailRedisRepository{redis}
}

func (c *cocktailRedisRepository) GetCocktailCollectionNumberLock(ctx context.Context, key string,
	ttl, retryInterval time.Duration, retryTimes int) (*redislock.Lock, error) {

	locker := redislock.New(c.redis)

	backoff := redislock.LimitRetry(redislock.LinearBackoff(retryInterval), retryTimes)

	// Obtain lock with retry, key is cocktail id
	lock, err := locker.Obtain(ctx, key, ttl, &redislock.Options{
		RetryStrategy: backoff,
	})
	if err == redislock.ErrNotObtained {
		return nil, domain.ErrRedisLockNotObtained
	} else if err != nil {
		return nil, err
	}

	return lock, nil
}

func (c *cocktailRedisRepository) ReleaseCocktailCollectionNumberLock(ctx context.Context, lock *redislock.Lock) error {
	return lock.Release(ctx)
}
