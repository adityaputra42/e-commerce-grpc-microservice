package services

import (
	"context"
	"e-commerce-microservice/auth/internal/config"
	"e-commerce-microservice/auth/internal/db"
	"e-commerce-microservice/auth/internal/model"
	"e-commerce-microservice/auth/internal/pb"
	"e-commerce-microservice/auth/internal/repository"
	"e-commerce-microservice/auth/internal/token"
	"e-commerce-microservice/auth/internal/utils"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
	RenewSessionLogin(ctx context.Context, refreshToekn string) (string, error)
}

type AuthServiceImpl struct {
	tokenMaker token.Maker
	config     config.Configuration
	repo       repository.AuthRepository
}

// RenewSessionLogin implements AuthService.
func (a *AuthServiceImpl) RenewSessionLogin(ctx context.Context, refreshToekn string) (string, error) {
	panic("unimplemented")
}

// Login implements AuthService.
func (a *AuthServiceImpl) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var userAuth pb.LoginResponse
	db := db.GetConnection()

	err := db.Transaction(func(tx *gorm.DB) error {
		user, err := a.repo.FindAuthUserByEmail(ctx, req.Email)
		if err != nil {
			return fmt.Errorf("User not found")
		}
		ok, err := utils.CheckPasswordHash(req.Password, user.HashedPassword)

		if !ok || err != nil {
			return fmt.Errorf("Password didn't match!")
		}

		AccessToken, _, err := a.tokenMaker.CreateToken(user.ID, utils.UserRole, a.config.AccessTokenDuration, token.TokenTypeAccessToken)
		if err != nil {
			return fmt.Errorf("Failed to create access token")
		}
		refreshToken, refreshPayload, err := a.tokenMaker.CreateToken(user.ID, utils.UserRole, a.config.RefreshTokenDuration, token.TokenTypeRefreshToken)
		if err != nil {
			return fmt.Errorf("Failed to create refresh token")
		}
		session := model.AuthSessions{
			ID:           refreshPayload.ID.String(),
			RefreshToken: refreshToken,
			UserId:       user.ID,
			UserAgent:    "",
			ClientIp:     "",
			IsBlocked:    false,
			ExpiredAt:    refreshPayload.ExpiredAt,
		}

		_, err = a.repo.CreateAuthSession(ctx, &session)
		if err != nil {
			return fmt.Errorf("Failed to create auth session")
		}
		userAuth = pb.LoginResponse{
			User: &pb.User{
				Id:         user.ID,
				Email:      user.HashedPassword,
				Provider:   user.Provider,
				IsVerified: user.IsVerified,
				UpdatedAt:  timestamppb.New(user.UpdatedAt),
				CreatedAt:  timestamppb.New(user.CreatedAt),
			},
			AccessToken:  AccessToken,
			RefreshToken: refreshToken,
		}

		return nil
	})

	if err != nil {
		return &pb.LoginResponse{}, err
	}

	return &userAuth, nil

}

// Register implements AuthService.
func (a *AuthServiceImpl) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var userAuth pb.RegisterResponse

	db := db.GetConnection()

	err := db.Transaction(func(tx *gorm.DB) error {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return fmt.Errorf("Failed to hash password")
		}
		req := model.AuthUsers{
			Email:          req.Email,
			HashedPassword: hashedPassword,
			Provider:       "email",
		}

		user, err := a.repo.CreateAuthUser(ctx, tx, &req)
		if err != nil {
			return fmt.Errorf("Failed to hash password")
		}
		userAuth = pb.RegisterResponse{
			User: &pb.User{
				Id:         user.ID,
				Email:      user.HashedPassword,
				Provider:   user.Provider,
				IsVerified: user.IsVerified,
				UpdatedAt:  timestamppb.New(user.UpdatedAt),
				CreatedAt:  timestamppb.New(user.CreatedAt),
			},
		}
		return nil
	})
	if err != nil {
		return &pb.RegisterResponse{}, err
	}

	return &userAuth, nil

}

func NewAuthServiceImpl(repo repository.AuthRepository) AuthService {
	return &AuthServiceImpl{repo: repo}
}
