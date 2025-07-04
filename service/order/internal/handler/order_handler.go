package handler

import (
	"context"
	"e-commerce-microservice/order/internal/pb"
	"e-commerce-microservice/order/internal/services"
	"e-commerce-microservice/order/internal/utils"
	"e-commerce-microservice/order/internal/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
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
	violations := validateCreateOrder(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := o.service.CreateOrder(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}

	return res, nil
}

func (o *OrderHandler) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.OrderResponse, error) {
	violations := validateId(req.GetId())
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := o.service.CancelOrder(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func (o *OrderHandler) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {

	violations := validateId(req.GetId())
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := o.service.DeleteOrder(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func (o *OrderHandler) FindAllOrder(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	violations := validateAllOrder(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := o.service.FindAllOrder(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func (o *OrderHandler) FindOrderById(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	violations := validateId(req.GetId())
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := o.service.FindOrderById(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func (o *OrderHandler) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error) {
	violations := validateId(req.GetId())
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := o.service.UpdateOrder(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func validateCreateOrder(req *pb.CreateOrderRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, utils.FieldViolation("username", err))
	}
	if err := val.ValidateUUID(req.GetCarId()); err != nil {
		violations = append(violations, utils.FieldViolation("car_id", err))
	}

	return violations
}

func validateId(id string) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUUID(id); err != nil {
		violations = append(violations, utils.FieldViolation("id", err))
	}

	return violations
}

func validateAllOrder(req *pb.ListOrdersRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, utils.FieldViolation("username", err))
	}

	return violations
}
