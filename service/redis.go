package service

import (
	"errors"
	"github.com/go-redis/redis"
	"time"
)

func newRedis(configure *Configure) (*redis.Client, error) {
	if configure.Redis == nil {
		return nil, errors.New("redis configure is not initialed")
	}
	redisConf := configure.Redis
	rdb := redis.NewClient(&redis.Options{
		Network:      redisConf.Network,
		Addr:      redisConf.Address,
		Password:     redisConf.Password,
		DB:           redisConf.DB,
		DialTimeout:  time.Second * time.Duration(redisConf.DialTimeoutSecond),
		ReadTimeout:  time.Second * time.Duration(redisConf.ReadTimeoutSecond),
		WriteTimeout: time.Second * time.Duration(redisConf.WriteTimeoutSecond),
		PoolSize:     redisConf.PoolSize,
	})

	return rdb, nil
}