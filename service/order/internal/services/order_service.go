package services

import (
	"context"
	db "e-commerce-microservice/order/internal/db/sqlc"
	"e-commerce-microservice/order/internal/pb"
	"e-commerce-microservice/order/internal/repository"
	"e-commerce-microservice/order/internal/token"
	"e-commerce-microservice/order/internal/utils"
	"e-commerce-microservice/order/internal/worker"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrdeService interface {
	CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error)
	UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error)
	CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.OrderResponse, error)
	DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error)
	FindOrderById(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error)
	FindAllOrder(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error)
}

type OrdeServiceImpl struct {
	Store      db.Store
	repo       repository.OrderRepository
	tokenMaker token.Maker
	natsWorker *worker.NatsWorker
}

// CancelOrder implements OrdeService.
func (o *OrdeServiceImpl) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.OrderResponse, error) {
	var orderUpdate pb.OrderResponse
	err := o.Store.ExecTx(ctx, func(q *db.Queries) error {

		_, err := utils.AuthorizationUser(ctx, o.tokenMaker)

		if err != nil {
			return utils.UnauthenticatedError(err)
		}

		order, err := o.repo.GetOrder(ctx, req.Id)

		if err != nil {
			return fmt.Errorf("order not found")
		}
		request := db.UpdateOrderParams{
			ID: order.ID,
			Status: pgtype.Text{
				String: utils.Candeled,
				Valid:  true,
			},
		}

		newOrder, err := o.repo.UpdateOrder(ctx, q, request)
		if err != nil {
			return fmt.Errorf("failed to cancel order")
		}
		orderUpdate = pb.OrderResponse{
			Order: &pb.Order{
				Id:        newOrder.ID.String(),
				Username:  newOrder.Username,
				CarId:     newOrder.CarID.String(),
				Status:    newOrder.Status,
				CreatedAt: timestamppb.New(newOrder.CreatedAt),
				UpdatedAt: timestamppb.New(newOrder.UpdatedAt),
			},
		}

		return nil
	})
	if err != nil {

		return nil, err
	}
	return &orderUpdate, nil
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
		amount, err := strconv.ParseFloat(req.GetAmount(), 64)
		if err != nil {
			return fmt.Errorf("failed to parse amount :%s", err)
		}
		param := db.CreateOrderParams{
			Username: req.GetUsername(),
			CarID:    carId,
			Amount:   amount,
		}
		result, err := o.repo.CreateOrder(ctx, tx, param)
		if err != nil {
			return fmt.Errorf("failed to create order: %s", err)
		}

		// Publish order created event
		if o.natsWorker != nil {
			err = o.natsWorker.PublishOrderCreated(ctx, worker.OrderCreatedPayload{
				OrderID:     result.ID.String(),
				UserID:      result.Username,
				TotalAmount: result.Amount,
				CreatedAt:   result.CreatedAt,
			})
			if err != nil {
				// Log the error but don't fail the order creation
				fmt.Printf("Failed to publish order created event: %v\n", err)
			}
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

	_, err := utils.AuthorizationUser(ctx, o.tokenMaker)

	if err != nil {
		return &pb.DeleteOrderResponse{Message: utils.UnauthenticatedError(err).Error()}, utils.UnauthenticatedError(err)
	}

	err = o.repo.DeleteOrder(ctx, req.Id)

	if err != nil {
		return &pb.DeleteOrderResponse{Message: "failed to delete order"}, fmt.Errorf("failed to delete order")
	}

	return &pb.DeleteOrderResponse{Message: "Success delete order"}, nil

}

// FindAllOrder implements OrdeService.
func (o *OrdeServiceImpl) FindAllOrder(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {

	var response pb.ListOrdersResponse

	authPayload, err := utils.AuthorizationUser(ctx, o.tokenMaker)

	if err != nil {
		return &pb.ListOrdersResponse{}, utils.UnauthenticatedError(err)
	}
	if authPayload.Username != req.Username {
		return &pb.ListOrdersResponse{}, fmt.Errorf("user not permitted")
	}
	request := db.ListOrderParams{
		Limit:  req.PageSize,
		Offset: req.Page,
	}
	orders, err := o.repo.GetListOrder(ctx, request)

	if err != nil {
		return &pb.ListOrdersResponse{}, fmt.Errorf("failed to get user orders")
	}

	orderPb := []*pb.Order{}

	for _, v := range *orders {
		orderPb = append(orderPb, &pb.Order{
			Id:        v.ID.String(),
			Username:  v.Username,
			CarId:     v.CarID.String(),
			Status:    v.Status,
			CreatedAt: timestamppb.New(v.CreatedAt),
			UpdatedAt: timestamppb.New(v.UpdatedAt),
		})
	}

	response = pb.ListOrdersResponse{
		Orders: orderPb,
	}

	return &response, nil
}

// FindOrderById implements OrdeService.
func (o *OrdeServiceImpl) FindOrderById(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	var order pb.OrderResponse

	_, err := utils.AuthorizationUser(ctx, o.tokenMaker)

	if err != nil {
		return nil, utils.UnauthenticatedError(err)
	}

	result, err := o.repo.GetOrder(ctx, req.GetId())

	if err != nil {
		return &pb.OrderResponse{}, fmt.Errorf("failed to get order")
	}

	order = pb.OrderResponse{
		Order: &pb.Order{
			Id:        result.ID.String(),
			Username:  result.Username,
			CarId:     result.CarID.String(),
			Status:    result.Status,
			CreatedAt: timestamppb.New(result.CreatedAt),
			UpdatedAt: timestamppb.New(result.UpdatedAt),
		},
	}

	return &order, nil
}

// UpdateOrder implements OrdeService.
func (o *OrdeServiceImpl) UpdateOrder(ctx context.Context, req *pb.UpdateOrderRequest) (*pb.OrderResponse, error) {
	var orderUpdate pb.OrderResponse
	err := o.Store.ExecTx(ctx, func(q *db.Queries) error {

		_, err := utils.AuthorizationUser(ctx, o.tokenMaker)

		if err != nil {
			return utils.UnauthenticatedError(err)
		}

		order, err := o.repo.GetOrder(ctx, req.Id)

		if err != nil {
			return fmt.Errorf("order not found")
		}
		request := db.UpdateOrderParams{
			ID: order.ID,
			Status: pgtype.Text{
				String: req.Status,
				Valid:  true,
			},
		}

		newOrder, err := o.repo.UpdateOrder(ctx, q, request)
		if err != nil {
			return fmt.Errorf("failed to update order")
		}
		orderUpdate = pb.OrderResponse{
			Order: &pb.Order{
				Id:        newOrder.ID.String(),
				Username:  newOrder.Username,
				CarId:     newOrder.CarID.String(),
				Status:    newOrder.Status,
				CreatedAt: timestamppb.New(newOrder.CreatedAt),
				UpdatedAt: timestamppb.New(newOrder.UpdatedAt),
			},
		}

		return nil
	})
	if err != nil {

		return nil, err
	}
	return &orderUpdate, nil
}

func NewOrderService(
	store db.Store,
	repo repository.OrderRepository,
	tokenMaker token.Maker,
	natsURL string) (OrdeService, error) {

	natsWorker, err := worker.NewNatsWorker(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create NATS worker: %v", err)
	}

	return &OrdeServiceImpl{
		Store:      store,
		repo:       repo,
		tokenMaker: tokenMaker,
		natsWorker: natsWorker,
	}, nil
}
