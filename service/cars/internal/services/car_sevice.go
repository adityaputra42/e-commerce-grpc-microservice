package services

import (
	"context"
	"e-commerce-microservice/cars/internal/model"
	"e-commerce-microservice/cars/internal/pb"
	"e-commerce-microservice/cars/internal/repository"
	"e-commerce-microservice/cars/internal/token"
	"e-commerce-microservice/cars/internal/utils"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type CarService interface {
	CreateCar(ctx context.Context, req *pb.CreateCarRequest) (*pb.CarResponse, error)
	FindCarById(ctx context.Context, req *pb.GetCarRequest) (*pb.CarResponse, error)
	UpdateCar(ctx context.Context, req *pb.UpdateCarRequest) (*pb.CarResponse, error)
	FindAllCars(ctx context.Context, req *pb.ListCarsRequest) (*pb.ListCarsResponse, error)
	DeletCar(ctx context.Context, req *pb.DeleteCarRequest) (*pb.DeleteCarResponse, error)
}
type CarServiceImpl struct {
	tokenMaker token.Maker
	db         *gorm.DB
	repo       repository.CarRepository
}

// CreateCar implements CarService.
func (c *CarServiceImpl) CreateCar(ctx context.Context, req *pb.CreateCarRequest) (*pb.CarResponse, error) {
	var repsonse pb.CarResponse
	err := c.db.Transaction(func(tx *gorm.DB) error {
		authPayload, err := utils.AuthorizationUser(ctx, c.tokenMaker)

		if err != nil {
			return utils.UnauthenticatedError(err)
		}
		if authPayload.Role != utils.AdminRole {
			return fmt.Errorf("user not permited")
		}
		carReq := model.Cars{
			Name:         req.GetName(),
			Brand:        req.GetBrand(),
			Model:        req.GetModel(),
			Year:         req.GetYear(),
			Mileage:      req.GetMileage(),
			Transmission: req.GetTransmission(),
			FuelType:     req.GetFuelType(),
			Location:     req.GetLocation(),
			Description:  req.GetDescription(),
			Images:       req.GetImages(),
			Price:        req.GetPrice(),
			Currency:     req.GetCurrency(),
		}
		result, err := c.repo.CreateCar(ctx, tx, &carReq)
		if err != nil {
			return fmt.Errorf("Failed to create Car")
		}
		repsonse = pb.CarResponse{
			Car: &pb.Car{
				Id:           result.ID,
				Name:         result.Name,
				Model:        result.Model,
				Brand:        result.Brand,
				Year:         result.Year,
				Mileage:      result.Mileage,
				Transmission: result.Transmission,
				FuelType:     result.FuelType,
				Location:     result.Location,
				Description:  result.Description,
				Images:       result.Images,
				Price:        result.Price,
				Currency:     result.Currency,
				IsSold:       result.IsSold,
				CreatedAt:    timestamppb.New(result.CreatedAt),
				UpdatedAt:    timestamppb.New(result.UpdatedAt),
			},
		}
		return nil
	})

	if err != nil {
		return &pb.CarResponse{}, err
	}

	return &repsonse, nil
}

// DeletCar implements CarService.
func (c *CarServiceImpl) DeletCar(ctx context.Context, req *pb.DeleteCarRequest) (*pb.DeleteCarResponse, error) {
	panic("unimplemented")
}

// FindAllCars implements CarService.
func (c *CarServiceImpl) FindAllCars(ctx context.Context, req *pb.ListCarsRequest) (*pb.ListCarsResponse, error) {
	panic("unimplemented")
}

// FindCarById implements CarService.
func (c *CarServiceImpl) FindCarById(ctx context.Context, req *pb.GetCarRequest) (*pb.CarResponse, error) {
	panic("unimplemented")
}

// UpdateCar implements CarService.
func (c *CarServiceImpl) UpdateCar(ctx context.Context, req *pb.UpdateCarRequest) (*pb.CarResponse, error) {
	panic("unimplemented")
}

func NewCarService(tokenMaker token.Maker,
	db *gorm.DB,
	repo repository.CarRepository) CarService {
	return &CarServiceImpl{tokenMaker: tokenMaker, db: db, repo: repo}
}
