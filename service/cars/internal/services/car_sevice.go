package services

import (
	"bytes"
	"context"
	"e-commerce-microservice/cars/internal/model"
	"e-commerce-microservice/cars/internal/pb"
	"e-commerce-microservice/cars/internal/repository"
	"e-commerce-microservice/cars/internal/token"
	"e-commerce-microservice/cars/internal/utils"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type CarService interface {
	CreateCar(ctx context.Context, req *pb.CreateCarRequest) (*pb.CarResponse, error)
	CreateCarWithImage(ctx context.Context, req *pb.CreateCarWithImageRequest) (*pb.CarResponse, error)
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
			return fmt.Errorf("failed to create Car")
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
	authPayload, err := utils.AuthorizationUser(ctx, c.tokenMaker)

	if err != nil {
		return &pb.DeleteCarResponse{Message: utils.UnauthenticatedError(err).Error()}, utils.UnauthenticatedError(err)
	}
	if authPayload.Role != utils.AdminRole {
		return &pb.DeleteCarResponse{Message: "user not permited"}, fmt.Errorf("user not permited")
	}

	result, err := c.repo.FindCarById(ctx, req.GetId())

	if err != nil {
		return &pb.DeleteCarResponse{Message: "failed to get car"}, fmt.Errorf("failed to get car")
	}
	err = c.repo.DeletCar(ctx, result)
	if err != nil {
		return &pb.DeleteCarResponse{Message: "failed to delete car"}, fmt.Errorf("failed to delete car")
	}
	return &pb.DeleteCarResponse{Message: "Success"}, nil

}

// CreateCarWithImage implements CarService.
func (c *CarServiceImpl) CreateCarWithImage(ctx context.Context, req *pb.CreateCarWithImageRequest) (*pb.CarResponse, error) {
	var repsonse pb.CarResponse
	err := c.db.Transaction(func(tx *gorm.DB) error {
		authPayload, err := utils.AuthorizationUser(ctx, c.tokenMaker)

		if err != nil {
			return utils.UnauthenticatedError(err)
		}
		if authPayload.Role != utils.AdminRole {
			return fmt.Errorf("user not permited")
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("image", req.GetImage().GetFilename())
		if err != nil {
			return err
		}
		_, err = io.Copy(part, bytes.NewReader(req.GetImage().GetContent()))
		if err != nil {
			return err
		}
		err = writer.Close()
		if err != nil {
			return err
		}

		r, err := http.NewRequest("POST", "http://localhost:8080/upload", body)
		if err != nil {
			return err
		}
		r.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, err := client.Do(r)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to upload image: %s", resp.Status)
		}

		var uploadResp struct {
			Urls []string `json:"urls"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
			return err
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
			Images:       []string{},
			Price:        req.GetPrice(),
			Currency:     req.GetCurrency(),
		}
		if len(uploadResp.Urls) > 0 {
			carReq.Images = []string{uploadResp.Urls[0]}
		}
		result, err := c.repo.CreateCar(ctx, tx, &carReq)
		if err != nil {
			return fmt.Errorf("failed to create Car")
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


// FindAllCars implements CarService.
func (c *CarServiceImpl) FindAllCars(ctx context.Context, req *pb.ListCarsRequest) (*pb.ListCarsResponse, error) {
	var cars pb.ListCarsResponse

	result, err := c.repo.FindAllCars(ctx)
	if err != nil {
		return &pb.ListCarsResponse{}, fmt.Errorf("failed to get cars")

	}
	listCar := []*pb.Car{}

	for _, v := range result {
		listCar = append(listCar, &pb.Car{
			Id:           v.ID,
			Name:         v.Name,
			Model:        v.Model,
			Brand:        v.Brand,
			Year:         v.Year,
			Mileage:      v.Mileage,
			Transmission: v.Transmission,
			FuelType:     v.FuelType,
			Location:     v.Location,
			Description:  v.Description,
			Images:       v.Images,
			Price:        v.Price,
			Currency:     v.Currency,
			IsSold:       v.IsSold,
			CreatedAt:    timestamppb.New(v.CreatedAt),
			UpdatedAt:    timestamppb.New(v.UpdatedAt),
		})
	}
	cars = pb.ListCarsResponse{
		Cars: listCar,
	}
	return &cars, nil

}

// FindCarById implements CarService.
func (c *CarServiceImpl) FindCarById(ctx context.Context, req *pb.GetCarRequest) (*pb.CarResponse, error) {
	var car pb.CarResponse

	result, err := c.repo.FindCarById(ctx, req.GetId())

	if err != nil {
		return nil, fmt.Errorf("failed to get car")
	}
	car = pb.CarResponse{
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

	return &car, nil

}

// UpdateCar implements CarService.
func (c *CarServiceImpl) UpdateCar(ctx context.Context, req *pb.UpdateCarRequest) (*pb.CarResponse, error) {
	var car pb.CarResponse

	err := c.db.Transaction(func(tx *gorm.DB) error {

		_, err := utils.AuthorizationUser(ctx, c.tokenMaker)

		if err != nil {
			return utils.UnauthenticatedError(err)
		}

		carUpdate, err := c.repo.FindCarById(ctx, req.GetId())

		if err != nil {
			return fmt.Errorf("failed to get car with id %s", req.Id)
		}

		request := model.Cars{
			ID:           carUpdate.ID,
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

		result, err := c.repo.UpdateCar(ctx, tx, &request)

		if err != nil {
			return fmt.Errorf("failed to update car")
		}
		car = pb.CarResponse{
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

	return &car, nil
}

func NewCarService(tokenMaker token.Maker,
	db *gorm.DB,
	repo repository.CarRepository) CarService {
	return &CarServiceImpl{tokenMaker: tokenMaker, db: db, repo: repo}
}
