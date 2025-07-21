package services

import (
	"context"
	"e-commerce-microservice/auth/internal/config"
	"e-commerce-microservice/auth/internal/model"
	pb "e-commerce-microservice/auth/internal/pb"
	"e-commerce-microservice/auth/internal/token"
	"e-commerce-microservice/auth/internal/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockAuthRepository is a mock for AuthRepository
type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) FindAuthUserByEmail(ctx context.Context, tx *gorm.DB, email string) (model.AuthUsers, error) {
	args := m.Called(ctx, tx, email)
	return args.Get(0).(model.AuthUsers), args.Error(1)
}

func (m *MockAuthRepository) CreateAuthUser(ctx context.Context, tx *gorm.DB, user *model.AuthUsers) (model.AuthUsers, error) {
	args := m.Called(ctx, tx, user)
	return args.Get(0).(model.AuthUsers), args.Error(1)
}

func (m *MockAuthRepository) UpdateAuthUser(ctx context.Context, tx *gorm.DB, user *model.AuthUsers) (model.AuthUsers, error) {
	args := m.Called(ctx, tx, user)
	return args.Get(0).(model.AuthUsers), args.Error(1)
}

func (m *MockAuthRepository) FindAuthUserById(ctx context.Context, username string) (model.AuthUsers, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(model.AuthUsers), args.Error(1)
}

func (m *MockAuthRepository) FindAuthSessionsById(ctx context.Context, id string) (model.AuthSessions, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.AuthSessions), args.Error(1)
}

func (m *MockAuthRepository) CreateAuthSession(ctx context.Context, tx *gorm.DB, session *model.AuthSessions) (model.AuthSessions, error) {
	args := m.Called(ctx, tx, session)
	return args.Get(0).(model.AuthSessions), args.Error(1)
}

func (m *MockAuthRepository) CreateVerifyEmail(ctx context.Context, verifyEmail *model.VerifyEmail) (model.VerifyEmail, error) {
	args := m.Called(ctx, verifyEmail)
	return args.Get(0).(model.VerifyEmail), args.Error(1)
}

func (m *MockAuthRepository) UpdateVerifyEmail(ctx context.Context, tx *gorm.DB, verifyEmail *model.VerifyEmail) (model.VerifyEmail, error) {
	args := m.Called(ctx, tx, verifyEmail)
	return args.Get(0).(model.VerifyEmail), args.Error(1)
}

func (m *MockAuthRepository) FindVerifyEmailById(ctx context.Context, id int64) (model.VerifyEmail, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.VerifyEmail), args.Error(1)
}

func TestLogin(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	cfg := config.Configuration{
		AccessTokenDuration:  time.Minute * 15,
		RefreshTokenDuration: time.Hour * 24,
	}

	tokenMaker, err := token.NewJWTMaker("12345678901234567890123456789012")
	require.NoError(t, err)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	service := &AuthServiceImpl{
		tokenMaker: tokenMaker,
		config:     cfg,
		repo:       mockRepo,
		db:         db,
	}

	// Create a proper hashed password for testing
	password := "password123"
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name    string
		req     *pb.LoginRequest
		role    string
		mock    func()
		wantErr bool
	}{
		{
			name: "successful login",
			req: &pb.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			role: "user",
			mock: func() {
				mockRepo.On("FindAuthUserByEmail", mock.Anything, mock.Anything, "test@example.com").
					Return(model.AuthUsers{
						Username:       "testuser",
						Email:          "test@example.com",
						HashedPassword: hashedPassword,
						Role:           "user",
					}, nil)

				session := model.AuthSessions{
					ID:           "test-session-id",
					RefreshToken: "test-refresh-token",
					Username:     "testuser",
					UserAgent:    "test-user-agent",
					ClientIp:     "127.0.0.1",
					IsBlocked:    false,
					ExpiredAt:    time.Now().Add(24 * time.Hour),
				}
				mockRepo.On("CreateAuthSession", mock.Anything, mock.Anything, mock.AnythingOfType("*model.AuthSessions")).
					Return(session, nil)

				mockRepo.On("CreateAuthSession", mock.Anything, mock.Anything, mock.AnythingOfType("*model.AuthSessions")).
					Return(&model.AuthSessions{}, nil)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			req: &pb.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			role: "user",
			mock: func() {
				mockRepo.On("FindAuthUserByEmail", mock.Anything, mock.Anything, "nonexistent@example.com").
					Return(model.AuthUsers{}, gorm.ErrRecordNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			resp, err := service.Login(context.Background(), tt.req, tt.role)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotEmpty(t, resp.AccessToken)
			require.NotEmpty(t, resp.RefreshToken)
		})
	}
}

func TestRenewSessionLogin(t *testing.T) {
	mockRepo := new(MockAuthRepository)
	cfg := config.Configuration{
		AccessTokenDuration:  time.Minute * 15,
		RefreshTokenDuration: time.Hour * 24,
	}

	tokenMaker, err := token.NewJWTMaker("12345678901234567890123456789012")
	require.NoError(t, err)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	service := &AuthServiceImpl{
		tokenMaker: tokenMaker,
		config:     cfg,
		repo:       mockRepo,
		db:         db,
	}

	tests := []struct {
		name    string
		setup   func() (string, *token.Payload)
		mock    func(string, *token.Payload)
		wantErr bool
	}{
		{
			name: "successful token refresh",
			setup: func() (string, *token.Payload) {
				refreshToken, payload, err := tokenMaker.CreateToken("testuser", "user", cfg.RefreshTokenDuration, token.TokenTypeRefreshToken)
				require.NoError(t, err)
				return refreshToken, payload
			},
			mock: func(refreshToken string, payload *token.Payload) {
				mockRepo.On("FindAuthSessionsById", mock.Anything, payload.ID.String()).
					Return(model.AuthSessions{
						ID:           payload.ID.String(),
						RefreshToken: refreshToken,
						Username:     "testuser",
						IsBlocked:    false,
						ExpiredAt:    time.Now().Add(time.Hour * 24),
					}, nil)
			},
			wantErr: false,
		},
		{
			name: "expired session",
			setup: func() (string, *token.Payload) {
				refreshToken, payload, err := tokenMaker.CreateToken("testuser", "user", time.Hour*-1, token.TokenTypeRefreshToken)
				require.NoError(t, err)
				return refreshToken, payload
			},
			mock: func(refreshToken string, payload *token.Payload) {
				mockRepo.On("FindAuthSessionsById", mock.Anything, payload.ID.String()).
					Return(model.AuthSessions{
						ID:           payload.ID.String(),
						RefreshToken: refreshToken,
						Username:     "testuser",
						IsBlocked:    false,
						ExpiredAt:    time.Now().Add(-time.Hour),
					}, nil)
					
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			refreshToken, payload := tt.setup()
			tt.mock(refreshToken, payload)

			resp, err := service.RenewSessionLogin(context.Background(), &pb.RefreshTokenRequest{
				RefreshToken: refreshToken,
			})

			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.NotEmpty(t, resp.AccessToken)
			require.Equal(t, refreshToken, resp.RefreshToken)
		})
	}
}
