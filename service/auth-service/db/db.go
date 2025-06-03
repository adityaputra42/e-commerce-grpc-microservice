package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB(DbSource string) {
	var err error
	dsn := DbSource
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic("Can't connect Database")
	}
	fmt.Println("Connected database")
}

func GetConnection() *gorm.DB {
	return db
}
