package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(DbSource string) (*gorm.DB, error) {
	var err error
	dsn := DbSource
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		return nil, fmt.Errorf("Can't connect Database")
	}
	fmt.Println("Connected database")
	return db, nil
}
