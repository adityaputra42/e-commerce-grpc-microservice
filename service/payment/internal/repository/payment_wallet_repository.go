package repository

import (
	"context"
	db "e-commerce-microservice/payment/internal/db/sqlc"
	"fmt"
)

type PaymentWalletRepository interface {
	FindAllPaymentWallet(ctx context.Context) (*[]db.PaymentWallet, error)
	CreatePaymentWallet(ctx context.Context, tx *db.Queries, req db.CreatePaymentWalletParams) (*db.PaymentWallet, error)
	DeletePaymentWallet(ctx context.Context, network string) error
	UpdatePaymentWallet(ctx context.Context, tx *db.Queries, req db.UpdatePaymentWalletParams) (*db.PaymentWallet, error)
}

type PaymentWalletRepositoryImpl struct {
	q *db.Queries
}

// CreatePaymentWallet implements PaymentWalletRepository.
func (p *PaymentWalletRepositoryImpl) CreatePaymentWallet(ctx context.Context, tx *db.Queries, req db.CreatePaymentWalletParams) (*db.PaymentWallet, error) {

	result, err := tx.CreatePaymentWallet(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Failed to create payment wallet")
	}
	return &result, nil
}

// DeletePaymentWallet implements PaymentWalletRepository.
func (p *PaymentWalletRepositoryImpl) DeletePaymentWallet(ctx context.Context, network string) error {

	err := p.q.DeletePaymentWallet(ctx, network)
	if err != nil {
		return fmt.Errorf("Failed to delete order")
	}
	return nil
}

// FindAllPaymentWallet implements PaymentWalletRepository.
func (p *PaymentWalletRepositoryImpl) FindAllPaymentWallet(ctx context.Context) (*[]db.PaymentWallet, error) {
	result, err := p.q.ListPaymentWallets(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get List payment wallet")
	}
	return &result, nil
}

// UpdatePaymentWallet implements PaymentWalletRepository.
func (p *PaymentWalletRepositoryImpl) UpdatePaymentWallet(ctx context.Context, tx *db.Queries, req db.UpdatePaymentWalletParams) (*db.PaymentWallet, error) {
	result, err := tx.UpdatePaymentWallet(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("Failed to update order")
	}
	return &result, nil
}

func NewPaymentWalletRepository(q *db.Queries) PaymentWalletRepository {
	return &PaymentWalletRepositoryImpl{q: q}
}
