package service

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
)

func newLogger(configure *Configure) (*logrus.Logger, error) {
	if configure.HTTP == nil {
		return nil, errors.New("logger configure is not initialed")
	}

	loggerConf := configure.Logger

	var log = logrus.New()

	level, err := logrus.ParseLevel(loggerConf.Level)
	if err != nil {
		return nil, err
	}
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(level)

	return log, nil
}

func GetLoggerEntry(log *logrus.Logger, ip, user_id string) *logrus.Entry {
	return log.WithFields(logrus.Fields{"request_ip": ip, "user_id": user_id})
}