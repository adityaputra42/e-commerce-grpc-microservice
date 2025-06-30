package services

import (
	"context"
	"e-commerce-microservice/auth/internal/config"
	"e-commerce-microservice/auth/internal/model"
	pb "e-commerce-microservice/auth/internal/pb"
	"e-commerce-microservice/auth/internal/repository"
	"e-commerce-microservice/auth/internal/token"
	"e-commerce-microservice/auth/internal/utils"
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, req *pb.RegisterRequest, role string) (*pb.RegisterResponse, error)
	Login(ctx context.Context, req *pb.LoginRequest, role string) (*pb.LoginResponse, error)
	RenewSessionLogin(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error)
}

type AuthServiceImpl struct {
	tokenMaker token.Maker
	config     config.Configuration
	db         *gorm.DB
	repo       repository.AuthRepository
}

// RenewSessionLogin implements AuthService.
func (a *AuthServiceImpl) RenewSessionLogin(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	payload, err := a.tokenMaker.VerifyToken(req.GetRefreshToken(), token.TokenTypeRefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	session, err := a.repo.FindAuthSessionsById(ctx, payload.ID.String())
	if err != nil || session.IsBlocked || session.ExpiredAt.Before(time.Now()) {
		return nil, fmt.Errorf("session invalid or expired")
	}

	accessToken, _, err := a.tokenMaker.CreateToken(payload.Username, payload.Role, a.config.AccessTokenDuration, token.TokenTypeAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token")
	}

	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken,
	}, nil
}

// Login implements AuthService.
func (a *AuthServiceImpl) Login(ctx context.Context, req *pb.LoginRequest, role string) (*pb.LoginResponse, error) {
	var userAuth pb.LoginResponse

	err := a.db.Transaction(func(tx *gorm.DB) error {
		user, err := a.repo.FindAuthUserByEmail(ctx, tx, req.Email)
		if err != nil {
			return fmt.Errorf("User not found")
		}
		ok, err := utils.CheckPasswordHash(req.Password, user.HashedPassword)

		if !ok || err != nil {
			return fmt.Errorf("Password didn't match!")
		}
		mtdt := utils.ExtractMetadata(ctx)
		log.Printf("userAgent: %v", mtdt.UserAgent)
		log.Printf("ClientIp: %v", mtdt.ClientIP)

		AccessToken, _, err := a.tokenMaker.CreateToken(user.Username, role, a.config.AccessTokenDuration, token.TokenTypeAccessToken)
		if err != nil {
			return fmt.Errorf("Failed to create access token")
		}
		refreshToken, refreshPayload, err := a.tokenMaker.CreateToken(user.Username, role, a.config.RefreshTokenDuration, token.TokenTypeRefreshToken)
		if err != nil {
			return fmt.Errorf("Failed to create refresh token")
		}
		session := model.AuthSessions{
			ID:           refreshPayload.ID.String(),
			RefreshToken: refreshToken,
			Username:     user.Username,
			UserAgent:    mtdt.UserAgent,
			ClientIp:     mtdt.ClientIP,
			IsBlocked:    false,
			ExpiredAt:    refreshPayload.ExpiredAt,
		}

		_, err = a.repo.CreateAuthSession(ctx, tx, &session)
		if err != nil {
			return fmt.Errorf("Failed to create auth session")
		}
		userAuth = pb.LoginResponse{
			User: &pb.AuthUser{
				Username:   user.Username,
				FullName:   user.FullName,
				Email:      user.Email,
				Role:       user.Role,
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
func (a *AuthServiceImpl) Register(ctx context.Context, req *pb.RegisterRequest, role string) (*pb.RegisterResponse, error) {
	var userAuth pb.RegisterResponse

	err := a.db.Transaction(func(tx *gorm.DB) error {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return fmt.Errorf("Failed to hash password")
		}
		req := model.AuthUsers{
			Username:       req.Username,
			FullName:       req.FullName,
			Email:          req.Email,
			HashedPassword: hashedPassword,
		}

		user, err := a.repo.CreateAuthUser(ctx, tx, &req)
		if err != nil {
			return fmt.Errorf("Failed to hash password")
		}
		mtdt := utils.ExtractMetadata(ctx)
		log.Printf("userAgent: %v", mtdt.UserAgent)
		log.Printf("ClientIp: %v", mtdt.ClientIP)

		AccessToken, _, err := a.tokenMaker.CreateToken(user.Username, role, a.config.AccessTokenDuration, token.TokenTypeAccessToken)
		if err != nil {
			return fmt.Errorf("Failed to create access token")
		}
		refreshToken, refreshPayload, err := a.tokenMaker.CreateToken(user.Username, role, a.config.RefreshTokenDuration, token.TokenTypeRefreshToken)
		if err != nil {
			return fmt.Errorf("Failed to create refresh token")
		}
		session := model.AuthSessions{
			ID:           refreshPayload.ID.String(),
			RefreshToken: refreshToken,
			Username:     user.Username,
			UserAgent:    mtdt.UserAgent,
			ClientIp:     mtdt.ClientIP,
			IsBlocked:    false,
			ExpiredAt:    refreshPayload.ExpiredAt,
		}

		_, err = a.repo.CreateAuthSession(ctx, tx, &session)
		if err != nil {
			return fmt.Errorf("Failed to create auth session")
		}

		userAuth = pb.RegisterResponse{
			User: &pb.AuthUser{
				Username:   user.Username,
				FullName:   user.FullName,
				Email:      user.Email,
				Role:       user.Role,
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
		return &pb.RegisterResponse{}, err
	}

	return &userAuth, nil

}

func NewAuthServiceImpl(repo repository.AuthRepository,
	tokenMaker token.Maker,
	config config.Configuration, db *gorm.DB) AuthService {
	return &AuthServiceImpl{repo: repo, tokenMaker: tokenMaker, config: config, db: db}
}
