package repository

import (
	"context"

	"github.com/adityaputra42/e-commerce-microservice/auth-service/internal/db"
	"github.com/adityaputra42/e-commerce-microservice/auth-service/internal/model"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateAuthUser(ctx context.Context, user model.AuthUsers) (model.AuthUsers, error)
	UpdateAuthUser(ctx context.Context, user model.AuthUsers) (model.AuthUsers, error)
	FindAuthUserById(ctx context.Context, id string) (model.AuthUsers, error)
	CreateAuthSession(ctx context.Context, user model.AuthSessions) (model.AuthSessions, error)
	FindAuthSessionsById(ctx context.Context, id string) (model.AuthSessions, error)
	CreateVerifyEmail(ctx context.Context, user model.VerifyEmail) (model.VerifyEmail, error)
	UpdateVerifyEmail(ctx context.Context, user model.VerifyEmail) (model.VerifyEmail, error)
	FindVerifyEmailById(ctx context.Context, id int64) (model.VerifyEmail, error)
}

type AuthRepositoryImpl struct {
	db *gorm.DB
}

// CreateAuthSession implements AuthRepository.
func (a AuthRepositoryImpl) CreateAuthSession(ctx context.Context, user model.AuthSessions) (model.AuthSessions, error) {
	panic("unimplemented")
}

// CreateAuthUser implements AuthRepository.
func (a AuthRepositoryImpl) CreateAuthUser(ctx context.Context, user model.AuthUsers) (model.AuthUsers, error) {
	panic("unimplemented")
}

// CreateVerifyEmail implements AuthRepository.
func (a AuthRepositoryImpl) CreateVerifyEmail(ctx context.Context, user model.VerifyEmail) (model.VerifyEmail, error) {
	panic("unimplemented")
}

// FindAuthSessionsById implements AuthRepository.
func (a AuthRepositoryImpl) FindAuthSessionsById(ctx context.Context, id string) (model.AuthSessions, error) {
	panic("unimplemented")
}

// FindAuthUserById implements AuthRepository.
func (a AuthRepositoryImpl) FindAuthUserById(ctx context.Context, id string) (model.AuthUsers, error) {
	panic("unimplemented")
}

// FindVerifyEmailById implements AuthRepository.
func (a AuthRepositoryImpl) FindVerifyEmailById(ctx context.Context, id int64) (model.VerifyEmail, error) {
	panic("unimplemented")
}

// UpdateAuthUser implements AuthRepository.
func (a AuthRepositoryImpl) UpdateAuthUser(ctx context.Context, user model.AuthUsers) (model.AuthUsers, error) {
	panic("unimplemented")
}

// UpdateVerifyEmail implements AuthRepository.
func (a AuthRepositoryImpl) UpdateVerifyEmail(ctx context.Context, user model.VerifyEmail) (model.VerifyEmail, error) {
	panic("unimplemented")
}

func NewSessionRepository() AuthRepository {
	db := db.GetConnection()
	return AuthRepositoryImpl{db: db}
}
