package redis

import (
	"context"
	"encoding/json"
	"github.com/beecool-cocktail/application-backend/domain"
	"github.com/go-redis/redis/v9"
)

type commandRedisRepository struct {
	redis *redis.Client
}

func NewRedisCommandRepository(redis *redis.Client) domain.CommandRedisRepository {
	return &commandRedisRepository{redis}
}

func (c *commandRedisRepository) Store(ctx context.Context, dc *domain.Command) error {
	key := "command:command_id:" + dc.ID
	value, err := json.Marshal(dc)
	if err != nil {
		return err
	}

	_, err = c.redis.Set(ctx, key, value, dc.ExpireTime).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *commandRedisRepository) GetByID(ctx context.Context, id string) (domain.Command, error) {

	var command domain.Command

	key := "command:command_id:" + id
	value, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		return domain.Command{}, err
	}

	err = json.Unmarshal([]byte(value), &command)
	if err != nil {
		return domain.Command{}, err
	}

	return command, err
}

func (c *commandRedisRepository) Delete(ctx context.Context, id string) error {
	key := "command:command_id:" + id

	_, err := c.redis.Del(ctx, key).Result()

	return err
}
