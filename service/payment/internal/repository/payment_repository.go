package repository

import (
	"context"
	db "e-commerce-microservice/payment/internal/db/sqlc"
	"fmt"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	FingPaymentById(ctx context.Context, id string) (*db.Payment, error)
	FindAllPayment(ctx context.Context, req db.ListPaymentsByUserParams) (*[]db.Payment, error)
	CreatePayment(ctx context.Context, tx *db.Queries, req db.CreatePaymentParams) (*db.Payment, error)

	UpdatePayment(ctx context.Context, tx *db.Queries, req db.UpdatePaymentStatusParams) (*db.Payment, error)
}
type PaymentRepositoryImpl struct {
	q *db.Queries
}

// CreatePayment implements PaymentRepository.
func (p *PaymentRepositoryImpl) CreatePayment(ctx context.Context, tx *db.Queries, req db.CreatePaymentParams) (*db.Payment, error) {

	result, err := tx.CreatePayment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Failed to create order")
	}
	return &result, nil
}

// FindAllPayment implements PaymentRepository.
func (p *PaymentRepositoryImpl) FindAllPayment(ctx context.Context, req db.ListPaymentsByUserParams) (*[]db.Payment, error) {
	result, err := p.q.ListPaymentsByUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get List Payment")
	}
	return &result, nil
}

// FingPaymentById implements PaymentRepository.
func (p *PaymentRepositoryImpl) FingPaymentById(ctx context.Context, id string) (*db.Payment, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid UUID: %s", err)
	}
	result, err := p.q.GetPayment(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get payment")
	}
	return &result, nil
}

// UpdatePayment implements PaymentRepository.
func (p *PaymentRepositoryImpl) UpdatePayment(ctx context.Context, tx *db.Queries, req db.UpdatePaymentStatusParams) (*db.Payment, error) {
	result, err := tx.UpdatePaymentStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Failed to update order")
	}
	return &result, nil
}

func NewPaymentRepository(q *db.Queries) PaymentRepository {
	return &PaymentRepositoryImpl{q: q}
}
