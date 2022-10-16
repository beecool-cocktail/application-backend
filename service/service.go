package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Service struct {
	DB        *gorm.DB
	Redis     *redis.Client
	Logger    *logrus.Logger
	HTTP      *gin.Engine
	Elastic   *elastic.Client
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

	elasticSearch, err := newElasticSearch(conf)
	if err != nil {
		return nil, err
	}

	service := &Service{
		Configure: conf,
		DB:        db,
		Redis:     rdb,
		HTTP:      http,
		Elastic:   elasticSearch,
		Logger:    logger,
	}

	return service, nil
}

func NewDBService(fileName string) (*Service, error) {
	conf, err := newConfigure(fileName)
	if err != nil {
		return nil, err
	}

	db, err := newMySQL(conf)
	if err != nil {
		return nil, err
	}

	service := &Service{
		Configure: conf,
		DB:        db,
	}

	return service, nil
}
