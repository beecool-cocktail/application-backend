package redis

import (
	"context"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/go-redis/redis"
	"strconv"
)

const cocktailCollectionNumbersKey = "cocktail_collection_numbers"

type cocktailRedisRepository struct {
	redis *redis.Client
}

func NewRedisCocktailRepository(redis *redis.Client) domain.CocktailRedisRepository {
	return &cocktailRedisRepository{redis}
}

func (c *cocktailRedisRepository) InitialCollectionNumbers(ctx context.Context, cr *domain.CocktailCollection) error {

	_, err := c.redis.ZAdd(
		cocktailCollectionNumbersKey,
		redis.Z{
			Score:  0,
			Member: cr.CocktailID,
		}).Result()

	return err
}

func (c *cocktailRedisRepository) IncreaseCollectionNumbers(ctx context.Context, cr *domain.CocktailCollection) error {
	_, err := c.redis.ZIncrBy(
		cocktailCollectionNumbersKey,
		cr.CollectionCounts,
		strconv.FormatInt(cr.CocktailID, 10),
	).Result()

	return err
}

func (c *cocktailRedisRepository) DecreaseCollectionNumbers(ctx context.Context, cr *domain.CocktailCollection) error {
	_, err := c.redis.ZIncrBy(
		cocktailCollectionNumbersKey,
		-cr.CollectionCounts,
		strconv.FormatInt(cr.CocktailID, 10),
	).Result()

	return err
}

func (c *cocktailRedisRepository) DeleteCollectionNumbers(ctx context.Context, cr *domain.CocktailCollection) error {
	_, err := c.redis.ZRem(
		cocktailCollectionNumbersKey,
		cr.CocktailID,
	).Result()

	return err
}
