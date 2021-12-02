package internal

import (
	"fmt"
	"github.com/beecool-cocktail/application-backend/service"
	"github.com/sirupsen/logrus"
)

func Migrate(cfgFile string, action string, tables []interface{}) {

	dbService, err := service.NewDBService(cfgFile)
	if err != nil {
		logrus.Panicf("start dbService failed - %s", err)
	}

	var n int

	//var err error
	switch action {
	case "migrate":
		count := 0
		for _, tbl := range tables {
			err = dbService.DB.AutoMigrate(tbl)
			if err != nil {
				panic(err)
			}
			count ++
		}
		n = count
	case "create":
		count := 0
		for _, tbl := range tables {
			err = dbService.DB.Migrator().CreateTable(tbl)
			if err != nil {
				panic(err)
			}
			count ++
		}
		n = count
	case "drop":
		count := 0
		for _, tbl := range tables {
			err = dbService.DB.Migrator().DropTable(tbl)
			if err != nil {
				panic(err)
			}
			count ++
		}
		n = count
	case "status":
		panic("status 尚未實作")
	case "import":
		panic("import 尚未實作")
	default:
		panic("action not found")
	}
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Applied %d migrations!\n", n)
}