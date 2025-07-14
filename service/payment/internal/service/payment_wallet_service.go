package service

import (
	"context"
	db "e-commerce-microservice/payment/internal/db/sqlc"
	"e-commerce-microservice/payment/internal/pb"
	"e-commerce-microservice/payment/internal/repository"
	"e-commerce-microservice/payment/internal/token"
	"e-commerce-microservice/payment/internal/utils"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PaymentWalletService interface {
	DeletePaymentWallet(ctx context.Context, req *pb.DeletePaymentWalletRequest) (*pb.DeletePaymentWalletResponse, error)
	FindPaymentWallet(ctx context.Context, req *pb.GetPaymentWalletRequest) (*pb.PaymentWalletResponse, error)
	FindAllPaymentWallet(ctx context.Context, req *pb.ListPaymentWalletsRequest) (*pb.ListPaymentWalletsResponse, error)
	CreatePaymentWalletParams(ctx context.Context, req *pb.CreatePaymentWalletRequest) (*pb.PaymentWalletResponse, error)
	UpdatePaymentWallet(ctx context.Context, req *pb.UpdatePaymentWalletRequest) (*pb.PaymentWalletResponse, error)
}
type PaymentWalletServiceImpl struct {
	Store      db.Store
	repo       repository.PaymentWalletRepository
	tokenMaker token.Maker
}

// FindPaymentWallet implements PaymentWalletService.
func (p *PaymentWalletServiceImpl) FindPaymentWallet(ctx context.Context, req *pb.GetPaymentWalletRequest) (*pb.PaymentWalletResponse, error) {
	var response pb.PaymentWalletResponse

	_, err := utils.AuthorizationUser(ctx, p.tokenMaker)

	if err != nil {
		return nil, utils.UnauthenticatedError(err)
	}

	result, err := p.repo.FindPaymentWalletByNetwork(ctx, req.GetNetwork())

	if err != nil {
		return &pb.PaymentWalletResponse{}, fmt.Errorf("failed to get payment wallet")
	}

	response = pb.PaymentWalletResponse{
		Wallet: &pb.PaymentWallet{
			Id:            result.ID,
			Network:       result.Network,
			WalletAddress: result.WalletAddress,
			CreatedAt:     timestamppb.New(result.CreatedAt.Time),
		},
	}

	return &response, nil
}

// CreatePaymentWalletParams implements PaymentWalletService.
func (p *PaymentWalletServiceImpl) CreatePaymentWalletParams(ctx context.Context, req *pb.CreatePaymentWalletRequest) (*pb.PaymentWalletResponse, error) {
	var response pb.PaymentWalletResponse
	err := p.Store.ExecTx(ctx, func(tx *db.Queries) error {
		_, err := utils.AuthorizationUser(ctx, p.tokenMaker)

		if err != nil {
			return utils.UnauthenticatedError(err)
		}

		param := db.CreatePaymentWalletParams{
			Network:       req.Network,
			WalletAddress: req.WalletAddress,
		}
		result, err := p.repo.CreatePaymentWallet(ctx, tx, param)
		if err != nil {
			return fmt.Errorf("failed to create payment: %s", err)
		}

		response = pb.PaymentWalletResponse{
			Wallet: &pb.PaymentWallet{
				Id:            result.ID,
				Network:       result.Network,
				WalletAddress: result.WalletAddress,
				CreatedAt:     timestamppb.New(result.CreatedAt.Time),
			},
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// DeletePaymentWallet implements PaymentWalletService.
func (p *PaymentWalletServiceImpl) DeletePaymentWallet(ctx context.Context, req *pb.DeletePaymentWalletRequest) (*pb.DeletePaymentWalletResponse, error) {
	authPayload, err := utils.AuthorizationUser(ctx, p.tokenMaker)

	if err != nil {
		return &pb.DeletePaymentWalletResponse{Message: utils.UnauthenticatedError(err).Error()}, utils.UnauthenticatedError(err)
	}
	if authPayload.Role != utils.AdminRole {
		return &pb.DeletePaymentWalletResponse{Message: "user not permited"}, fmt.Errorf("user not permited")
	}

	network, err := p.repo.FindPaymentWalletByNetwork(ctx, req.GetNetwork())
	if err != nil {
		return &pb.DeletePaymentWalletResponse{Message: "payment wallet not found"}, fmt.Errorf("payment wallet not found")

	}
	err = p.repo.DeletePaymentWallet(ctx, network.Network)
	if err != nil {
		return &pb.DeletePaymentWalletResponse{Message: "failed delete payment wallet"}, fmt.Errorf("failed delete payment wallet")

	}
	return &pb.DeletePaymentWalletResponse{Message: "success"}, nil

}

// FindAllPaymentWallet implements PaymentWalletService.
func (p *PaymentWalletServiceImpl) FindAllPaymentWallet(ctx context.Context, req *pb.ListPaymentWalletsRequest) (*pb.ListPaymentWalletsResponse, error) {

	var response pb.ListPaymentWalletsResponse

	_, err := utils.AuthorizationUser(ctx, p.tokenMaker)

	if err != nil {
		return &pb.ListPaymentWalletsResponse{}, utils.UnauthenticatedError(err)
	}

	payments, err := p.repo.FindAllPaymentWallet(ctx)

	if err != nil {
		return &pb.ListPaymentWalletsResponse{}, fmt.Errorf("failed to get payment wallet list")
	}

	paymentPb := []*pb.PaymentWallet{}

	for _, v := range *payments {
		paymentPb = append(paymentPb, &pb.PaymentWallet{
			Id:            v.ID,
			Network:       v.Network,
			WalletAddress: v.WalletAddress,
			CreatedAt:     timestamppb.New(v.CreatedAt.Time),
		})
	}

	response = pb.ListPaymentWalletsResponse{
		Wallets: paymentPb,
	}

	return &response, nil
}

// UpdatePaymentWallet implements PaymentWalletService.
func (p *PaymentWalletServiceImpl) UpdatePaymentWallet(ctx context.Context, req *pb.UpdatePaymentWalletRequest) (*pb.PaymentWalletResponse, error) {
	var response pb.PaymentWalletResponse
	err := p.Store.ExecTx(ctx, func(q *db.Queries) error {

		_, err := utils.AuthorizationUser(ctx, p.tokenMaker)

		if err != nil {
			return utils.UnauthenticatedError(err)
		}

		payment, err := p.repo.FindPaymentWalletByNetwork(ctx, req.Network)

		if err != nil {
			return fmt.Errorf("payment wallet not found")
		}
		request := db.UpdatePaymentWalletParams{
			Network:       payment.Network,
			WalletAddress: req.WalletAddress,
		}

		result, err := p.repo.UpdatePaymentWallet(ctx, q, request)
		if err != nil {
			return fmt.Errorf("failed to update order")
		}

		response = pb.PaymentWalletResponse{
			Wallet: &pb.PaymentWallet{
				Id:            result.ID,
				Network:       result.Network,
				WalletAddress: result.WalletAddress,
				CreatedAt:     timestamppb.New(result.CreatedAt.Time),
			},
		}

		return nil
	})
	if err != nil {

		return nil, err
	}
	return &response, nil
}

func NewPayementWalletService(Store db.Store,
	repo repository.PaymentWalletRepository,
	tokenMaker token.Maker) PaymentWalletService {
	return &PaymentWalletServiceImpl{Store: Store, repo: repo, tokenMaker: tokenMaker}
}
