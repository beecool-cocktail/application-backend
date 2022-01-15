package service

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func newLogger(configure *Configure) (*logrus.Logger, error) {
	if configure.Logger == nil {
		return nil, errors.New("logger configure is not initialed")
	}

	loggerConf := configure.Logger

	var log = logrus.New()

	level, err := logrus.ParseLevel(loggerConf.Level)
	if err != nil {
		return nil, err
	}
	//log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(level)

	return log, nil
}

func GetLoggerEntry(log *logrus.Logger, api string, request interface{}) *logrus.Entry {
	formatRequest := fmt.Sprintf("%+v", request)
	return log.WithFields(logrus.Fields{"api": api, "request_data": formatRequest})

}