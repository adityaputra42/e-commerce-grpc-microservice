package handler

import (
	"context"
	"e-commerce-microservice/user/internal/pb"
	"e-commerce-microservice/user/internal/services"
	"e-commerce-microservice/user/internal/utils"
	"e-commerce-microservice/user/internal/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser implements UserHandler
func (a *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	violations := validateCreateUser(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}

	res, err := a.service.CreateUser(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

// CreateUserAddress implements UserHandler
func (a *UserHandler) CreateUserAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.UserAddressResponse, error) {
	violations := validateCreateAddress(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}

	res, err := a.service.CreateUserAddress(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

// FindAddressUser implements UserHandler
func (a *UserHandler) FindAllAddress(ctx context.Context, req *pb.GetUserAddressesRequest) (*pb.UserAddressesResponse, error) {
	violations := validateGetAddress(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}

	res, err := a.service.FindUserAddresses(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

// UpdateUserAddress implements UserHandler
func (a *UserHandler) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.UserAddressResponse, error) {
	violations := validateUpdateAddress(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}
	res, err := a.service.UpdateUserAddresses(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func validateCreateUser(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, utils.FieldViolation("username", err))
	}

	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, utils.FieldViolation("full_name", err))
	}
	if err := val.ValidatePhoneNumber(req.GetPhoneNumber()); err != nil {
		violations = append(violations, utils.FieldViolation("phone_number", err))
	}
	return violations
}

func validateCreateAddress(req *pb.CreateAddressRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, utils.FieldViolation("username", err))
	}

	return violations
}

func validateGetAddress(req *pb.GetUserAddressesRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, utils.FieldViolation("username", err))
	}

	return violations
}

func validateUpdateAddress(req *pb.UpdateAddressRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUUID(req.GetId()); err != nil {
		violations = append(violations, utils.FieldViolation("id", err))
	}

	return violations
}
