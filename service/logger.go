package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type Logger struct {
	Logger   *logrus.Logger
	Path     string
	FileName string
}

func newLogger(configure *Configure) (*Logger, error) {
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

	if _, err := os.Stat(configure.Logger.Path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(configure.Logger.Path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	logger := &Logger{
		Logger:   log,
		Path:     loggerConf.Path,
		FileName: loggerConf.FileName,
	}

	return logger, nil
}

func (l *Logger) GetLoggerFields(user int64, ip, method string, requestBody interface{}, uri interface{}) logrus.Fields {

	formatRequest := fmt.Sprintf("%+v", requestBody)
	return logrus.Fields{"user": user, "ip": ip, "method": method, "uri": uri, "requestBody": formatRequest}
}

func (l *Logger) LogFile(ctx context.Context, level logrus.Level, fields logrus.Fields, format string, args ...interface{}) {

	filePath := l.Path + l.FileName + "-" + time.Now().Format("20060102") + ".log"
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	l.Logger.SetOutput(os.Stdout)
	l.Logger.SetOutput(f)

	switch level {
	case logrus.PanicLevel:
		if len(args) > 0 {
			l.Logger.WithFields(fields).Panicf(format, args...)
		} else {
			l.Logger.WithFields(fields).Panic(format)
		}
	case logrus.FatalLevel:
		if len(args) > 0 {
			l.Logger.WithFields(fields).Fatalf(format, args...)
		} else {
			l.Logger.WithFields(fields).Fatal(format)
		}
	case logrus.ErrorLevel:
		if len(args) > 0 {
			l.Logger.WithFields(fields).Errorf(format, args...)
		} else {
			l.Logger.WithFields(fields).Error(format)
		}
	case logrus.WarnLevel:
		if len(args) > 0 {
			l.Logger.WithFields(fields).Warnf(format, args...)
		} else {
			l.Logger.WithFields(fields).Warn(format)
		}
	case logrus.InfoLevel:
		if len(args) > 0 {
			l.Logger.WithFields(fields).Infof(format, args...)
		} else {
			l.Logger.WithFields(fields).Info(format)
		}
	case logrus.DebugLevel:
		if len(args) > 0 {
			l.Logger.WithFields(fields).Debugf(format, args...)
		} else {
			l.Logger.WithFields(fields).Debug(format)
		}
	default:
		if len(args) > 0 {
			l.Logger.WithFields(fields).Infof(format, args...)
		} else {
			l.Logger.WithFields(fields).Info(format)
		}
	}
	return
}
