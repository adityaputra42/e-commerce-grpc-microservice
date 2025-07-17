package handler

import (
	"context"
	"e-commerce-microservice/cars/internal/pb"
	"e-commerce-microservice/cars/internal/services"
	"e-commerce-microservice/cars/internal/utils"
	"e-commerce-microservice/cars/internal/val"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CarHandler struct {
	pb.UnimplementedCarServiceServer
	service services.CarService
}

func NewCarHandler(service services.CarService) *CarHandler {
	return &CarHandler{service: service}
}

func (c *CarHandler) CreateCar(ctx context.Context, req *pb.CreateCarRequest) (*pb.CarResponse, error) {
	violations := validateCreateCar(req)
	if violations != nil {
		return nil, utils.InvalidArgumentError(violations)

	}

	response, err := c.service.CreateCar(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return response, nil
}

func (c *CarHandler) CreateCarWithImage(ctx context.Context, req *pb.CreateCarWithImageRequest) (*pb.CarResponse, error) {
	// no validation for now
	response, err := c.service.CreateCarWithImage(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return response, nil
}


func (a *CarHandler) FindAllCars(ctx context.Context, req *pb.ListCarsRequest) (*pb.ListCarsResponse, error) {

	res, err := a.service.FindAllCars(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func (a *CarHandler) UpdateCar(ctx context.Context, req *pb.UpdateCarRequest) (*pb.CarResponse, error) {
	res, err := a.service.UpdateCar(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func (a *CarHandler) DeleteCar(ctx context.Context, req *pb.DeleteCarRequest) (*pb.DeleteCarResponse, error) {
	res, err := a.service.DeletCar(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func (a *CarHandler) FindCarById(ctx context.Context, req *pb.GetCarRequest) (*pb.CarResponse, error) {
	res, err := a.service.FindCarById(ctx, req)
	if err != nil {
		status.Errorf(codes.Internal, "%s", err.Error())
	}
	return res, nil
}

func validateCreateCar(req *pb.CreateCarRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateName(req.GetName()); err != nil {
		violations = append(violations, utils.FieldViolation("name", err))
	}

	if err := val.ValidateBrand(req.GetBrand()); err != nil {
		violations = append(violations, utils.FieldViolation("brand", err))
	}

	if err := val.ValidateModel(req.GetModel()); err != nil {
		violations = append(violations, utils.FieldViolation("model", err))
	}

	if err := val.ValidateYear(strconv.Itoa(int(req.GetYear()))); err != nil {
		violations = append(violations, utils.FieldViolation("year", err))
	}

	if err := val.ValidateMileage(strconv.Itoa(int(req.GetMileage()))); err != nil {
		violations = append(violations, utils.FieldViolation("mileage", err))
	}
	if err := val.ValidateTransmission(req.GetTransmission()); err != nil {
		violations = append(violations, utils.FieldViolation("transmission", err))
	}
	if err := val.ValidateFullType(req.GetFuelType()); err != nil {
		violations = append(violations, utils.FieldViolation("fuel_type", err))
	}
	if err := val.ValidateLocation(req.GetLocation()); err != nil {
		violations = append(violations, utils.FieldViolation("location", err))
	}
	if err := val.ValidateDescription(req.GetDescription()); err != nil {
		violations = append(violations, utils.FieldViolation("description", err))
	}
	if err := val.ValidateCurrency(req.GetCurrency()); err != nil {
		violations = append(violations, utils.FieldViolation("currency", err))
	}
	return violations
}
