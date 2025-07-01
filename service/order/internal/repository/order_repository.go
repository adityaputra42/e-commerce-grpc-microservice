package repository

import (
	"context"
	db "e-commerce-microservice/order/internal/db/sqlc"
	"fmt"

	"github.com/google/uuid"
)

type OrderRepository interface {
	GetOrder(ctx context.Context, id string) (*db.Order, error)
	GetListOrder(ctx context.Context, req db.ListOrderParams) (*[]db.Order, error)
	CreateOrder(ctx context.Context, req db.CreateOrderParams) (*db.Order, error)
	DeleteOrder(ctx context.Context, id string) error
	UpdateOrder(ctx context.Context, req db.UpdateOrderParams) (*db.Order, error)
}

type OrderRepositoryImpl struct {
	q *db.Queries
}

// CreateOrder implements OrderRepository.
func (o *OrderRepositoryImpl) CreateOrder(ctx context.Context, req db.CreateOrderParams) (*db.Order, error) {

	result, err := o.q.CreateOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Failed to create order")
	}
	return &result, nil
}

// DeleteOrder implements OrderRepository.
func (o *OrderRepositoryImpl) DeleteOrder(ctx context.Context, id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid UUID: %s", err)
	}
	err = o.q.DeleteOrder(ctx, uuid)
	if err != nil {
		return fmt.Errorf("Failed to delete order")
	}
	return nil
}

// GetListOrder implements OrderRepository.
func (o *OrderRepositoryImpl) GetListOrder(ctx context.Context, req db.ListOrderParams) (*[]db.Order, error) {
	result, err := o.q.ListOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get List order")
	}
	return &result, nil
}

// GetOrder implements OrderRepository.
func (o *OrderRepositoryImpl) GetOrder(ctx context.Context, id string) (*db.Order, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %s", err)
	}
	result, err := o.q.GetOrder(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get order")
	}
	return &result, nil
}

// UpdaeOrder implements OrderRepository.
func (o *OrderRepositoryImpl) UpdateOrder(ctx context.Context, req db.UpdateOrderParams) (*db.Order, error) {
	result, err := o.q.UpdateOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Failed to update order")
	}
	return &result, nil
}

func NewOrderRepository(q *db.Queries) OrderRepository {
	return &OrderRepositoryImpl{q: q}
}
