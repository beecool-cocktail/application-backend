package testutil

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	UserID  int64  = 1
	Account string = "Andy"
)

func GetLogger() *logrus.Entry {
	log := logrus.New()
	logger := logrus.NewEntry(log)
	return logger
}

func GetRouter() *gin.Engine {
	r := gin.Default()

	return r
}

func GetRouteWithcontext() *gin.Engine {
	r := gin.Default()
	r.Use(setAdminInfo())

	return r
}

func setAdminInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", UserID)
		c.Set("account", Account)
	}
}

func BeforeEach() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}