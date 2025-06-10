package repository

import (
	"context"
	"e-commerce-microservice/user/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, tx *gorm.DB, user *model.Users) (model.Users, error)
	FindUserById(ctx context.Context, username string) (model.Users, error)
	UpdateUser(ctx context.Context, tx *gorm.DB, user *model.Users) (model.Users, error)
	CreateUserAddress(ctx context.Context, tx *gorm.DB, userAddress *model.UserAddresses) (model.UserAddresses, error)
	FindUserAdrressById(ctx context.Context, id string) (model.UserAddresses, error)
	FindUserAdrressesByUsername(ctx context.Context, username string) ([]model.UserAddresses, error)
	UpdateUserAddress(ctx context.Context, tx *gorm.DB, userAddress *model.UserAddresses) (model.UserAddresses, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

// CreateUser implements UserRepository.
func (u *UserRepositoryImpl) CreateUser(ctx context.Context, tx *gorm.DB, user *model.Users) (model.Users, error) {
	result := tx.WithContext(ctx).Create(&user)
	if result.Error != nil {
		return model.Users{}, result.Error
	}

	return *user, nil
}

// CreateUserAddress implements UserRepository.
func (u *UserRepositoryImpl) CreateUserAddress(ctx context.Context, tx *gorm.DB, userAddress *model.UserAddresses) (model.UserAddresses, error) {
	result := tx.WithContext(ctx).Create(&userAddress)
	if result.Error != nil {
		return model.UserAddresses{}, result.Error
	}

	return *userAddress, nil
}

// FindUserAdrressById implements UserRepository.
func (u *UserRepositoryImpl) FindUserAdrressById(ctx context.Context, id string) (model.UserAddresses, error) {
	userAddress := model.UserAddresses{}
	err := u.db.WithContext(ctx).Model(&model.UserAddresses{}).Take(&userAddress, "id =?", id).Error

	if err != nil {
		return model.UserAddresses{}, err
	}
	return userAddress, nil
}

// FindUserAdrressesByUsername implements UserRepository.
func (u *UserRepositoryImpl) FindUserAdrressesByUsername(ctx context.Context, username string) ([]model.UserAddresses, error) {
	userAddress := []model.UserAddresses{}
	err := u.db.WithContext(ctx).
		Where("username = ?", username).
		Find(&userAddress).Error
	if err != nil {
		return []model.UserAddresses{}, err
	}
	return userAddress, nil
}

// FindUserById implements UserRepository.
func (u *UserRepositoryImpl) FindUserById(ctx context.Context, username string) (model.Users, error) {
	user := model.Users{}
	err := u.db.WithContext(ctx).Model(&model.Users{}).Take(&user, "username =?", username).Error

	if err != nil {
		return model.Users{}, err
	}
	return user, nil
}

// UpdateUser implements UserRepository.
func (u *UserRepositoryImpl) UpdateUser(ctx context.Context, tx *gorm.DB, user *model.Users) (model.Users, error) {
	result := tx.WithContext(ctx).Save(&user)

	if result.Error != nil {
		return model.Users{}, result.Error
	}

	return *user, nil
}

// UpdateUserAddress implements UserRepository.
func (u *UserRepositoryImpl) UpdateUserAddress(ctx context.Context, tx *gorm.DB, userAddress *model.UserAddresses) (model.UserAddresses, error) {
	result := tx.WithContext(ctx).Save(&userAddress)

	if result.Error != nil {
		return model.UserAddresses{}, result.Error
	}

	return *userAddress, nil
}

func NewUserRepositoryImpl(DB *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: DB}
}
