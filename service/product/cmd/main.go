package main

import (
	"context"
	"e-commerce-microservice/product/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("Cannot load config")
	}

	_, err = pgxpool.New(context.Background(), config.DbSource)

	if err != nil {
		log.Fatal().Msg("cannot connect to db:")
	}
}
