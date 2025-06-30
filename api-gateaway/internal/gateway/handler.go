package gateway

import (
	"context"
	"e-commerce-microservice/api-gateway/internal/pb"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewGatewayMux(ctx context.Context) http.Handler {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		log.Fatalf("cannot register auth service: %v", err)
	}

	if err := pb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, "localhost:50052", opts); err != nil {
		log.Fatalf("cannot register user service: %v", err)
	}

	if err := pb.RegisterCarServiceHandlerFromEndpoint(ctx, mux, "localhost:50053", opts); err != nil {
		log.Fatalf("cannot register car service: %v", err)
	}

	return mux
}
