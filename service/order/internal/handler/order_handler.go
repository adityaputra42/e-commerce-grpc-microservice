package handler

import (
	"context"
	"e-commerce-microservice/order/internal/pb"
	"e-commerce-microservice/order/internal/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	pb.UnimplementedOrderServiceServer
	service services.OrdeService
}

func NewOrderHandler(service services.OrdeService) *OrderHandler {
	return &OrderHandler{service: service}
}
func (o *OrderHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	res, err := o.service.CreateOrder(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}

	return res, nil
}

func (o *OrderHandler) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.OrderResponse, error) {

	res, err := o.service.CancelOrder(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func (o *OrderHandler) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	res, err := o.service.DeleteOrder(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func (o *OrderHandler) FindAllOrder(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {

	res, err := o.service.FindAllOrder(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func (o *OrderHandler) FindOrderById(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {

	res, err := o.service.FindOrderById(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func (o *OrderHandler) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error) {

	res, err := o.service.UpdateOrder(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}
