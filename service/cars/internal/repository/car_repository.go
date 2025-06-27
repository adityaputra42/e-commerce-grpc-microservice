package repository

import (
	"context"
	"e-commerce-microservice/cars/internal/model"

	"gorm.io/gorm"
)

type CarRepository interface {
	CreateCar(ctx context.Context, tx *gorm.DB, cars *model.Cars) (model.Cars, error)
	UpdateCar(ctx context.Context, tx *gorm.DB, cars *model.Cars) (model.Cars, error)
	FindCarById(ctx context.Context, id string) (model.Cars, error)
	FindAllCars(ctx context.Context) ([]model.Cars, error)
	DeletCar(ctx context.Context, cars model.Cars) error
}

type CarRepositoryImpl struct {
	db *gorm.DB
}

// CreateCar implements CarRepository.
func (c *CarRepositoryImpl) CreateCar(ctx context.Context, tx *gorm.DB, cars *model.Cars) (model.Cars, error) {
	result := tx.WithContext(ctx).Create(&cars)
	if result.Error != nil {
		return model.Cars{}, result.Error
	}

	return *cars, nil
}

// DeletCar implements CarRepository.
func (c *CarRepositoryImpl) DeletCar(ctx context.Context, cars model.Cars) error {
	result := c.db.Delete(&cars)
	return result.Error
}

// FindAllCars implements CarRepository.
func (c *CarRepositoryImpl) FindAllCars(ctx context.Context) ([]model.Cars, error) {
	cars := []model.Cars{}
	err := c.db.WithContext(ctx).
		Find(&cars).Error
	if err != nil {
		return []model.Cars{}, err
	}
	return cars, nil
}

// FindCarById implements CarRepository.
func (c *CarRepositoryImpl) FindCarById(ctx context.Context, id string) (model.Cars, error) {
	car := model.Cars{}
	err := c.db.WithContext(ctx).Model(&model.Cars{}).Take(&car, "id =?", id).Error

	if err != nil {
		return model.Cars{}, err
	}
	return car, nil
}

// UpdateCar implements CarRepository.
func (c *CarRepositoryImpl) UpdateCar(ctx context.Context, tx *gorm.DB, cars *model.Cars) (model.Cars, error) {
	result := tx.WithContext(ctx).Save(&cars)
	if result.Error != nil {
		return model.Cars{}, result.Error
	}
	return *cars, nil
}

func NewCarRepository(db *gorm.DB) CarRepository {
	return &CarRepositoryImpl{}
}
