package cmd

import (
	"e-commerce-microservice/cars/internal/config"
	"e-commerce-microservice/cars/internal/db"

	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
		panic(err)
	}

	_, err = db.InitDB(conf.DbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize database")
		panic(err)
	}
}
