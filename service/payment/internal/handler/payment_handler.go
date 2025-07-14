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

type PaymentHandler struct {
	pb.UnimplementedPaymentServiceServer
	service service.PaymentService
}

// CreatePayment implements PaymentInterface.
func (p *PaymentHandler) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.PaymentResponse, error) {
	violations := validateCreatePayment(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := p.service.CreatePayment(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}

	return res, nil
}

// FindAllPayment implements PaymentInterface.
func (p *PaymentHandler) FindAllPayment(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	violations := validateGetAll(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := p.service.FindAllPayment(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

// FingPaymentById implements PaymentInterface.
func (p *PaymentHandler) FingPaymentById(ctx context.Context, req *pb.GetPaymentRequest) (*pb.PaymentResponse, error) {
	violations := validateGetById(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := p.service.FingPaymentById(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

// UpdatePayment implements PaymentInterface.
func (p *PaymentHandler) UpdatePayment(ctx context.Context, req *pb.UpdatePaymentStatusRequest) (*pb.PaymentResponse, error) {
	violations := validateUpdatePayement(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := p.service.UpdatePayment(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}

	return res, nil
}

func NewPaymentHandler(Service service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: Service}
}

func validateCreatePayment(req *pb.CreatePaymentRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUUID(req.GetOrderId()); err != nil {
		violations = append(violations, utils.FieldViolation("order_id", err))
	}
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, utils.FieldViolation("username", err))
	}
	if err := val.ValidateNetwork(req.GetNetwork()); err != nil {
		violations = append(violations, utils.FieldViolation("network", err))
	}
	if err := val.ValidateCurrency(req.GetCurrency()); err != nil {
		violations = append(violations, utils.FieldViolation("currency", err))
	}
	if err := val.ValidateAmount(req.GetAmount()); err != nil {
		violations = append(violations, utils.FieldViolation("amount", err))
	}
	if err := val.ValidateWalletAddress(req.GetNetwork(), req.GetWalletAddress()); err != nil {
		violations = append(violations, utils.FieldViolation("wallet_address", err))
	}

	return violations
}

func validateUpdatePayement(req *pb.UpdatePaymentStatusRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUUID(req.GetId()); err != nil {
		violations = append(violations, utils.FieldViolation("id", err))
	}
	if err := val.ValidateTxHash(req.GetNetwork(), req.GetTxHash()); err != nil {
		violations = append(violations, utils.FieldViolation("tx_hash", err))
	}

	return violations
}

func validateGetAll(req *pb.ListPaymentsRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, utils.FieldViolation("username", err))
	}

	return violations
}

func validateGetById(req *pb.GetPaymentRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUUID(req.GetId()); err != nil {
		violations = append(violations, utils.FieldViolation("id", err))
	}

	return violations
}
