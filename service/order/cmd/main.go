package main

import (
	"context"
	"e-commerce-microservice/order/internal/config"
	db "e-commerce-microservice/order/internal/db/sqlc"
	"e-commerce-microservice/order/internal/handler"
	"e-commerce-microservice/order/internal/pb"
	"e-commerce-microservice/order/internal/repository"
	"e-commerce-microservice/order/internal/services"
	"e-commerce-microservice/order/internal/token"
	"e-commerce-microservice/order/internal/utils"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("Cannot load config")
	}
	shutdownCtx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()
	waitGroup, ctx := errgroup.WithContext(shutdownCtx)

	connPool, err := pgxpool.New(ctx, config.DbSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db:")
	}
	store := db.NewStore(connPool)

	runGrpcServer(ctx, waitGroup, config, store)
	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}

}
func runGrpcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config config.Configuration,
	db *db.SQLStore,

) {

	tokenMaker, err := token.NewJWTMaker(config.SecretKey)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create token maker")

	}
	orderRepo := repository.NewOrderRepository(db.Queries)
	orderService := services.NewOrderService(db, orderRepo, tokenMaker)
	orderHandler := handler.NewOrderHandler(orderService)

	grpcLogger := grpc.UnaryInterceptor(utils.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterOrderServiceServer(grpcServer, orderHandler)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC Order server at %s", listener.Addr().String())

		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("gRPC Order server failed to serve")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown gRPC Ordr server")

		grpcServer.GracefulStop()
		log.Info().Msg("gRPC Order server is stopped")

		return nil
	})
}
