package main

import (
	"context"
	"e-commerce-microservice/cars/internal/config"
	"e-commerce-microservice/cars/internal/db"
	"e-commerce-microservice/cars/internal/handler"
	"e-commerce-microservice/cars/internal/pb"
	"e-commerce-microservice/cars/internal/repository"
	"e-commerce-microservice/cars/internal/services"
	"e-commerce-microservice/cars/internal/token"
	"e-commerce-microservice/cars/internal/utils"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
		panic(err)
	}
	shutdownCtx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()
	waitGroup, ctx := errgroup.WithContext(shutdownCtx)

	db, err := db.InitDB(conf.DbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize database")
		panic(err)
	}
	runGrpcServer(ctx, waitGroup, conf, db)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}
}

func runGrpcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config config.Configuration,
	db *gorm.DB,

) {

	tokenMaker, err := token.NewJWTMaker(config.SecretKey)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create token maker")

	}
	carRepo := repository.NewCarRepository(db)
	carService := services.NewCarService(tokenMaker, db, carRepo)
	carHandler := handler.NewCarHandler(carService)

	grpcLogger := grpc.UnaryInterceptor(utils.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterCarServiceServer(grpcServer, carHandler)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC cars server at %s", listener.Addr().String())

		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("gRPC cars server failed to serve")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown gRPC cars server")

		grpcServer.GracefulStop()
		log.Info().Msg("gRPC cars server is stopped")

		return nil
	})
}
