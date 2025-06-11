package main

import (
	"context"
	"e-commerce-microservice/user/internal/config"
	"e-commerce-microservice/user/internal/db"
	"e-commerce-microservice/user/internal/handler"
	"e-commerce-microservice/user/internal/pb"
	"e-commerce-microservice/user/internal/repository"
	"e-commerce-microservice/user/internal/services"
	"e-commerce-microservice/user/internal/token"
	"e-commerce-microservice/user/internal/utils"
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
	userRepo := repository.NewUserRepositoryImpl(db)
	userService := services.NewUserService(tokenMaker, db, userRepo)
	userHandler := handler.NewUserHandler(userService)

	grpcLogger := grpc.UnaryInterceptor(utils.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterUserServiceServer(grpcServer, userHandler)
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
