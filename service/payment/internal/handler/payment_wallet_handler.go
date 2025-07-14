package handler

import (
	"context"
	"e-commerce-microservice/payment/internal/pb"
	"e-commerce-microservice/payment/internal/service"
	"e-commerce-microservice/payment/internal/utils"
	"e-commerce-microservice/payment/internal/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentWalletHandler struct {
	pb.UnimplementedPaymentWalletServiceServer
	service service.PaymentWalletService
}

// CreatePaymentWalletParams implements PaymentWallett.
func (p *PaymentWalletHandler) CreatePaymentWalletParams(ctx context.Context, req *pb.CreatePaymentWalletRequest) (*pb.PaymentWalletResponse, error) {
	violations := validateCreatePaymentWallet(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := p.service.CreatePaymentWallet(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}

	return res, nil
}

// DeletePaymentWallet implements PaymentWallett.
func (p *PaymentWalletHandler) DeletePaymentWallet(ctx context.Context, req *pb.DeletePaymentWalletRequest) (*pb.DeletePaymentWalletResponse, error) {
	violations := validateDeletePayementWallet(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := p.service.DeletePaymentWallet(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}

	return res, nil
}

// FindAllPaymentWallet implements PaymentWallett.
func (p *PaymentWalletHandler) FindAllPaymentWallet(ctx context.Context, req *pb.ListPaymentWalletsRequest) (*pb.ListPaymentWalletsResponse, error) {
	res, err := p.service.FindAllPaymentWallet(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}

	return res, nil
}

// FindPaymentWallet implements PaymentWallett.
func (p *PaymentWalletHandler) FindPaymentWallet(ctx context.Context, req *pb.GetPaymentWalletRequest) (*pb.PaymentWalletResponse, error) {
	violations := validateGetPayementWallet(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := p.service.FindPaymentWallet(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}

	return res, nil
}

// UpdatePaymentWallet implements PaymentWallett.
func (p *PaymentWalletHandler) UpdatePaymentWallet(ctx context.Context, req *pb.UpdatePaymentWalletRequest) (*pb.PaymentWalletResponse, error) {
	violations := validateUpdatePayementWallet(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := p.service.UpdatePaymentWallet(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}

	return res, nil
}

func NewPaymentWalletHandler(s service.PaymentWalletService) *PaymentWalletHandler {
	return &PaymentWalletHandler{service: s}
}

func validateCreatePaymentWallet(req *pb.CreatePaymentWalletRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateNetwork(req.GetNetwork()); err != nil {
		violations = append(violations, utils.FieldViolation("network", err))
	}

	if err := val.ValidateWalletAddress(req.GetNetwork(), req.GetWalletAddress()); err != nil {
		violations = append(violations, utils.FieldViolation("wallet_address", err))
	}

	return violations
}

func validateUpdatePayementWallet(req *pb.UpdatePaymentWalletRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateNetwork(req.GetNetwork()); err != nil {
		violations = append(violations, utils.FieldViolation("network", err))
	}
	if err := val.ValidateWalletAddress(req.GetNetwork(), req.GetWalletAddress()); err != nil {
		violations = append(violations, utils.FieldViolation("wallet_address", err))
	}

	return violations
}

func validateGetPayementWallet(req *pb.GetPaymentWalletRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateNetwork(req.GetNetwork()); err != nil {
		violations = append(violations, utils.FieldViolation("network", err))
	}

	return violations
}

func validateDeletePayementWallet(req *pb.DeletePaymentWalletRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateNetwork(req.GetNetwork()); err != nil {
		violations = append(violations, utils.FieldViolation("network", err))
	}

	return violations
}
