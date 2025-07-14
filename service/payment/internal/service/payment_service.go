package service

import (
	"context"
	db "e-commerce-microservice/payment/internal/db/sqlc"
	"e-commerce-microservice/payment/internal/pb"
	"e-commerce-microservice/payment/internal/repository"
	"e-commerce-microservice/payment/internal/token"
	"e-commerce-microservice/payment/internal/utils"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentService interface {
	FingPaymentById(ctx context.Context, req *pb.GetPaymentRequest) (*pb.PaymentResponse, error)
	FindAllPayment(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error)
	CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.PaymentResponse, error)
	UpdatePayment(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.PaymentResponse, error)
}
type PaymentServiceImpl struct {
	Store      db.Store
	repo       repository.PaymentRepository
	tokenMaker token.Maker
}

// CreatePayment implements PaymentService.
func (p *PaymentServiceImpl) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.PaymentResponse, error) {
	var response pb.PaymentResponse
	err := p.Store.ExecTx(ctx, func(tx *db.Queries) error {
		authPayload, err := utils.AuthorizationUser(ctx, p.tokenMaker)

		if err != nil {
			return utils.UnauthenticatedError(err)
		}
		if authPayload.Username != req.Username {
			return fmt.Errorf("username mismatch")
		}
		orderId, err := uuid.Parse(req.GetOrderId())
		if err != nil {
			return fmt.Errorf("invalid payment ID: %s", err)
		}

		amount, err := strconv.ParseFloat(req.Amount, 64)
		if err != nil {
			return fmt.Errorf("invalid Amount: %s", err)
		}
		param := db.CreatePaymentParams{
			OrderID:       orderId,
			Username:      req.Username,
			Network:       req.Network,
			Currency:      req.Currency,
			Amount:        amount,
			WalletAddress: req.WalletAddress,
		}
		result, err := p.repo.CreatePayment(ctx, tx, param)
		if err != nil {
			return fmt.Errorf("failed to create payment: %s", err)
		}

		response = pb.PaymentResponse{
			Payment: &pb.Payment{
				Id:        result.ID.String(),
				Username:  result.Username,
				OrderId:   result.OrderID.String(),
				Network:   result.Network,
				Currency:  result.Currency,
				Amount:    strconv.FormatFloat(result.Amount, 'f', -1, 64),
				TxHash:    result.TxHash.String,
				Status:    result.Status,
				CreatedAt: timestamppb.New(result.CreatedAt.Time),
				UpdatedAt: timestamppb.New(result.UpdatedAt.Time),
			},
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// FindAllPayment implements PaymentService.
func (p *PaymentServiceImpl) FindAllPayment(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {

	var response pb.ListPaymentsResponse

	authPayload, err := utils.AuthorizationUser(ctx, p.tokenMaker)

	if err != nil {
		return &pb.ListPaymentsResponse{}, utils.UnauthenticatedError(err)
	}
	if authPayload.Username != req.Username {
		return &pb.ListPaymentsResponse{}, fmt.Errorf("user not permitted")
	}
	request := db.ListPaymentsByUserParams{
		Username: req.Username,
		Limit:    req.PageSize,
		Offset:   req.Page,
	}
	payments, err := p.repo.FindAllPayment(ctx, request)

	if err != nil {
		return &pb.ListPaymentsResponse{}, fmt.Errorf("failed to get payment list")
	}

	paymentPb := []*pb.Payment{}

	for _, v := range *payments {
		paymentPb = append(paymentPb, &pb.Payment{
			Id:        v.ID.String(),
			Username:  v.Username,
			OrderId:   v.OrderID.String(),
			Network:   v.Network,
			Currency:  v.Currency,
			Amount:    strconv.FormatFloat(v.Amount, 'f', -1, 64),
			TxHash:    v.TxHash.String,
			Status:    v.Status,
			CreatedAt: timestamppb.New(v.CreatedAt.Time),
			UpdatedAt: timestamppb.New(v.UpdatedAt.Time),
		})
	}

	response = pb.ListPaymentsResponse{
		Payments: paymentPb,
	}

	return &response, nil
}

// FingPaymentById implements PaymentService.
func (p *PaymentServiceImpl) FingPaymentById(ctx context.Context, req *pb.GetPaymentRequest) (*pb.PaymentResponse, error) {
	var response pb.PaymentResponse

	_, err := utils.AuthorizationUser(ctx, p.tokenMaker)

	if err != nil {
		return nil, utils.UnauthenticatedError(err)
	}

	result, err := p.repo.FingPaymentById(ctx, req.GetId())

	if err != nil {
		return &pb.PaymentResponse{}, fmt.Errorf("failed to get payment")
	}

	response = pb.PaymentResponse{
		Payment: &pb.Payment{
			Id:        result.ID.String(),
			Username:  result.Username,
			OrderId:   result.OrderID.String(),
			Network:   result.Network,
			Currency:  result.Currency,
			Amount:    strconv.FormatFloat(result.Amount, 'f', -1, 64),
			TxHash:    result.TxHash.String,
			Status:    result.Status,
			CreatedAt: timestamppb.New(result.CreatedAt.Time),
			UpdatedAt: timestamppb.New(result.UpdatedAt.Time),
		},
	}

	return &response, nil
}

// UpdatePayment implements PaymentService.
func (p *PaymentServiceImpl) UpdatePayment(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.PaymentResponse, error) {
	var paymentUpdate pb.PaymentResponse
	err := p.Store.ExecTx(ctx, func(q *db.Queries) error {

		_, err := utils.AuthorizationUser(ctx, p.tokenMaker)

		if err != nil {
			return utils.UnauthenticatedError(err)
		}

		payment, err := p.repo.FingPaymentById(ctx, req.Id)

		if err != nil {
			return fmt.Errorf("payment not found")
		}
		request := db.UpdatePaymentStatusParams{
			ID:     payment.ID,
			Status: req.Status,
			TxHash: pgtype.Text{
				String: req.TxHash,
				Valid:  true,
			},
		}

		result, err := p.repo.UpdatePayment(ctx, q, request)
		if err != nil {
			return fmt.Errorf("failed to update order")
		}
		paymentUpdate = pb.PaymentResponse{
			Payment: &pb.Payment{
				Id:        result.ID.String(),
				Username:  result.Username,
				OrderId:   result.OrderID.String(),
				Network:   result.Network,
				Currency:  result.Currency,
				Amount:    strconv.FormatFloat(result.Amount, 'f', -1, 64),
				TxHash:    result.TxHash.String,
				Status:    result.Status,
				CreatedAt: timestamppb.New(result.CreatedAt.Time),
				UpdatedAt: timestamppb.New(result.UpdatedAt.Time),
			},
		}

		return nil
	})
	if err != nil {

		return nil, err
	}
	return &paymentUpdate, nil
}

func NewPayementService(Store db.Store,
	repo repository.PaymentRepository,
	tokenMaker token.Maker) PaymentService {
	return &PaymentServiceImpl{Store: Store, repo: repo, tokenMaker: tokenMaker}
}
