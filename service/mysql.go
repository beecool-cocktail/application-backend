package service

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func newMySQL(configure *Configure) (*gorm.DB, error) {
	if configure.DB == nil {
		return nil, errors.New("mysql configure is not initialed")
	}
	mainDBConf := configure.DB.MainDB

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True",
		mainDBConf.User, mainDBConf.Password, mainDBConf.Address, mainDBConf.DBName)

	dbConnection, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{DisableForeignKeyConstraintWhenMigrating:true})
	if err != nil {
		return nil, err
	}

	db, err := dbConnection.DB()
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(time.Duration(mainDBConf.SetConnMaxIdleTime) * time.Second)
	db.SetMaxIdleConns(mainDBConf.SetMaxIdleConns)
	db.SetMaxOpenConns(mainDBConf.SetMaxOpenConns)


	return dbConnection, nil
}