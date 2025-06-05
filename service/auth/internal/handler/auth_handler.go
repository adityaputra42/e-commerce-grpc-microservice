package handler

import (
	"context"
	"e-commerce-microservice/auth/internal/pb"
	"e-commerce-microservice/auth/internal/services"
	"e-commerce-microservice/auth/internal/utils"
	"e-commerce-microservice/auth/internal/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler interface {
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
}

type AuthHandlerImpl struct {
	pb.UnimplementedAuthServiceServer
	service services.AuthService
}

// Login implements AuthHandler.
func (a *AuthHandlerImpl) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	violations := validateLoginUser(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}

	res, err := a.service.Login(ctx, req, utils.UserRole)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

// Register implements AuthHandler.
func (a *AuthHandlerImpl) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	violations := validateRegisterUser(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}

	res, err := a.service.Register(ctx, req, utils.UserRole)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func NewAuthServiceImpl(authService services.AuthService) AuthHandler {
	return &AuthHandlerImpl{service: authService}
}

func validateLoginUser(req *pb.LoginRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, utils.FieldViolation("email", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, utils.FieldViolation("password", err))
	}
	return violations
}

func validateRegisterUser(req *pb.RegisterRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, utils.FieldViolation("email", err))
	}

	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, utils.FieldViolation("password", err))
	}
	return violations
}
