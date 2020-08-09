package ormutil

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var defaultDB *gorm.DB

func Init(driver, dsn string) {
	db, err := gorm.Open(driver, dsn)
	if err != nil {
		log.Fatalln(err)
	}

	defaultDB = db
}

func Close() {
	defaultDB.Close()
}

func DB() *gorm.DB {

	return defaultDB
}
