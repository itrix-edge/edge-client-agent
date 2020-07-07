package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
)

type ORM struct{}

var dbDSN string

// Usage: orm.DB(); orm. ...
var orm *gorm.DB

// var sqlDB *sql.DB
// func init() {
// 	o.Init()
// }

func Init(dbDSN string, debug bool) {
	var err error
	if orm, err = OpenConnection(dbDSN, debug); err != nil {
		log.Printf("failed to connect database, got error %v\n", err)
		os.Exit(1)
	} else {
		sqlDB, err := orm.DB()
		if err != nil {
			log.Printf("failed to get database, got error %v\n", err)
		}
		err = sqlDB.Ping()
		if err != nil {
			log.Printf("failed to connect database, got error %v\n", err)
		}
	}
	// orm, err := OpenConnection()
}

func OpenConnection(dbDSN string, debug bool) (db *gorm.DB, err error) {

	// dbDSN := fmt.Sprintf("user=%s password=%s DB.name=%s host=%s port=%s %s", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_OPTIONS"))
	fmt.Println(dbDSN)
	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dbDSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if debug {
		db.Logger = db.Logger.LogMode(logger.Info)
	} else {
		db.Logger = db.Logger.LogMode(logger.Silent)
	}

	return db, err
}

func GetORM() *gorm.DB {
	return orm
}
