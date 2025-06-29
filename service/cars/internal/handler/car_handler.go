package handler

import (
	"context"
	"e-commerce-microservice/cars/internal/pb"
	"e-commerce-microservice/cars/internal/services"

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
	response, err := c.service.CreateCar(ctx, req)
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
