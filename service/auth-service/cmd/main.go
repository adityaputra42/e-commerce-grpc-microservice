package main

import (
	"log"

	"github.com/adityaputra42/e-commerce-microservice/auth-service/config"
	"github.com/adityaputra42/e-commerce-microservice/auth-service/db"
)

func main() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config", err)
		panic(err)
	}

	db.InitDB(conf.DbSource)
}
