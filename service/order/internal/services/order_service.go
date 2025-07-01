package services

import (
	"context"
	db "e-commerce-microservice/order/internal/db/sqlc"
	"e-commerce-microservice/order/internal/pb"
	"e-commerce-microservice/order/internal/repository"
	"e-commerce-microservice/order/internal/token"
)

type OrdeService interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error)
	UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error)
	DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error)
	FindOrderById(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error)
	FindAllOrder(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error)
}

type name struct {
	store      *db.Store
	repo       repository.OrderRepository
	tokenMaker token.Maker
}
