package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Service struct {
	DB        *gorm.DB
	Redis     *redis.Client
	Logger    *logrus.Logger
	HTTP      *gin.Engine
	Configure *Configure
}

func NewService(fileName string) (*Service, error) {
	conf, err := newConfigure(fileName)
	if err != nil {
		return nil, err
	}

	db, err := newMySQL(conf)
	if err != nil {
		return nil, err
	}

	rdb, err := newRedis(conf)
	if err != nil {
		return nil, err
	}

	http, err := newHTTP(conf)
	if err != nil {
		return nil, err
	}

	logger, err := newLogger(conf)
	if err != nil {
		return nil, err
	}

	service := &Service{
		Configure: conf,
		DB:        db,
		Redis:     rdb,
		HTTP:      http,
		Logger:    logger,
	}

	return service, nil
}
