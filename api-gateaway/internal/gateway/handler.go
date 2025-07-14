package gateway

import (
	"context"
	"e-commerce-microservice/api-gateway/internal/config"
	"e-commerce-microservice/api-gateway/internal/pb"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGatewayMux(ctx context.Context, config config.Configuration) http.Handler {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, config.AuthServerAddress, opts); err != nil {
		log.Fatalf("cannot register auth service: %v", err)
	}

	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, config.UserServerAddress, opts); err != nil {
		log.Fatalf("cannot register user service: %v", err)
	}

	if err := pb.RegisterCarServiceHandlerFromEndpoint(ctx, mux, config.CarsServerAddress, opts); err != nil {
		log.Fatalf("cannot register car service: %v", err)
	}

	if err := pb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, config.OrderServerAddress, opts); err != nil {
		log.Fatalf("cannot register order service: %v", err)
	}
	if err := pb.RegisterPaymentServiceHandlerFromEndpoint(ctx, mux, config.PaymentServerAddress, opts); err != nil {
		log.Fatalf("cannot register payment service: %v", err)
	}
	if err := pb.RegisterPaymentWalletServiceHandlerFromEndpoint(ctx, mux, config.PaymentServerAddress, opts); err != nil {
		log.Fatalf("cannot register payment wallet service: %v", err)
	}

	return mux
}
