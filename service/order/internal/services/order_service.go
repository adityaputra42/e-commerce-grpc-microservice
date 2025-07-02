package services

import (
	"context"
	db "e-commerce-microservice/order/internal/db/sqlc"
	"e-commerce-microservice/order/internal/pb"
	"e-commerce-microservice/order/internal/repository"
	"e-commerce-microservice/order/internal/token"
	"e-commerce-microservice/order/internal/utils"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrdeService interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error)
	UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error)
	CancelOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error)
	DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error)
	FindOrderById(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error)
	FindAllOrder(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error)
}

type OrdeServiceImpl struct {
	Store      db.Store
	repo       repository.OrderRepository
	tokenMaker token.Maker
}

// CancelOrder implements OrdeService.
func (o *OrdeServiceImpl) CancelOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error) {
	panic("unimplemented")
}

// CreateOrder implements OrdeService.
func (o *OrdeServiceImpl) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	var response pb.OrderResponse
	err := o.Store.ExecTx(ctx, func(tx *db.Queries) error {
		authPayload, err := utils.AuthorizationUser(ctx, o.tokenMaker)

		if err != nil {
			return utils.UnauthenticatedError(err)
		}
		if authPayload.Username != req.Username {
			return fmt.Errorf("username mismatch")
		}
		carId, err := uuid.Parse(req.GetCarId())
		if err != nil {
			return fmt.Errorf("invalid Car UUID: %s", err)
		}
		param := db.CreateOrderParams{
			Username: req.GetUsername(),
			CarID:    carId,
		}
		result, err := o.repo.CreateOrder(ctx, tx, param)
		if err != nil {
			return fmt.Errorf("failed to create order: %s", err)
		}

		response = pb.OrderResponse{
			Order: &pb.Order{
				Id:        result.ID.String(),
				Username:  result.Username,
				CarId:     result.CarID.String(),
				Status:    result.Status,
				CreatedAt: timestamppb.New(result.CreatedAt),
				UpdatedAt: timestamppb.New(result.UpdatedAt),
			},
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// DeleteOrder implements OrdeService.
func (o *OrdeServiceImpl) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {
	panic("unimplemented")
}

// FindAllOrder implements OrdeService.
func (o *OrdeServiceImpl) FindAllOrder(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	panic("unimplemented")
}

// FindOrderById implements OrdeService.
func (o *OrdeServiceImpl) FindOrderById(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	panic("unimplemented")
}

// UpdateOrder implements OrdeService.
func (o *OrdeServiceImpl) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error) {
	panic("unimplemented")
}

func NewOrderService(
	store db.Store,
	repo repository.OrderRepository,
	tokenMaker token.Maker) OrdeService {
	return &OrdeServiceImpl{Store: store, repo: repo, tokenMaker: tokenMaker}
}
