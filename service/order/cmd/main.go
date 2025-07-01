package main

import (
	"context"
	"e-commerce-microservice/order/internal/config"
	db "e-commerce-microservice/order/internal/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("Cannot load config")
	}
	connPool, err := pgxpool.New(context.Background(), config.DbSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db:")
	}
	db.NewStore(connPool)

}
