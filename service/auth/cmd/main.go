package main

import (
	"e-commerce-microservice/auth/internal/config"
	"e-commerce-microservice/auth/internal/db"
	"log"
)

func main() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
		panic(err)
	}

	db.InitDB(conf.DbSource)
}
