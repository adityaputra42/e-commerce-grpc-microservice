package main

import (
	"context"
	"e-commerce-microservice/auth/internal/config"
	"e-commerce-microservice/auth/internal/db"
	"e-commerce-microservice/auth/internal/handler"
	"e-commerce-microservice/auth/internal/pb"
	"e-commerce-microservice/auth/internal/repository"
	"e-commerce-microservice/auth/internal/services"
	"e-commerce-microservice/auth/internal/token"
	"e-commerce-microservice/auth/internal/utils"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"net"

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
	conf, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
		panic(err)
	}
	shutdownCtx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()
	waitGroup, ctx := errgroup.WithContext(shutdownCtx)

	db.InitDB(conf.DbSource)
	runGrpcServer(ctx, waitGroup, conf)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("error from wait group")
	}

}

func runGrpcServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	config config.Configuration,

) {

	tokenMaker, err := token.NewJWTMaker(config.SecretKey)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create token maker")

	}
	userRepo := repository.NewAuthRepository()
	userService := services.NewAuthServiceImpl(userRepo, tokenMaker, config)
	userHandler := handler.NewAuthHandler(userService)

	grpcLogger := grpc.UnaryInterceptor(utils.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterAuthServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create listener")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msg("gRPC server failed to serve")
			return err
		}

		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown gRPC server")

		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server is stopped")

		return nil
	})
}
