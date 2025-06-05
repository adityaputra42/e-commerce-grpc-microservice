package repository

import (
	"context"
	"e-commerce-microservice/auth/internal/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateAuthUser(ctx context.Context, tx *gorm.DB, authUsers *model.AuthUsers) (model.AuthUsers, error)
	UpdateAuthUser(ctx context.Context, tx *gorm.DB, authUsers *model.AuthUsers) (model.AuthUsers, error)
	FindAuthUserById(ctx context.Context, username string) (model.AuthUsers, error)
	FindAuthUserByEmail(ctx context.Context, tx *gorm.DB, email string) (model.AuthUsers, error)
	CreateAuthSession(ctx context.Context, tx *gorm.DB, authSession *model.AuthSessions) (model.AuthSessions, error)
	FindAuthSessionsById(ctx context.Context, id string) (model.AuthSessions, error)
	CreateVerifyEmail(ctx context.Context, verifyEmail *model.VerifyEmail) (model.VerifyEmail, error)
	UpdateVerifyEmail(ctx context.Context, tx *gorm.DB, verifyEmail *model.VerifyEmail) (model.VerifyEmail, error)
	FindVerifyEmailById(ctx context.Context, id int64) (model.VerifyEmail, error)
}

type AuthRepositoryImpl struct {
	db *gorm.DB
}

// FindAuthUserEmail implements AuthRepository.
func (a AuthRepositoryImpl) FindAuthUserByEmail(ctx context.Context, tx *gorm.DB, email string) (model.AuthUsers, error) {
	authUser := model.AuthUsers{}
	err := a.db.WithContext(ctx).Model(&model.AuthUsers{}).Take(&authUser, "email =?", email).Error
	if err != nil {
		return model.AuthUsers{}, err
	}
	return authUser, nil
}

// CreateAuthSession implements AuthRepository.
func (a AuthRepositoryImpl) CreateAuthSession(ctx context.Context, tx *gorm.DB, authSession *model.AuthSessions) (model.AuthSessions, error) {
	result := a.db.WithContext(ctx).Create(&authSession)
	if result.Error != nil {
		return model.AuthSessions{}, result.Error
	}

	return *authSession, nil
}

// CreateAuthUser implements AuthRepository.
func (a AuthRepositoryImpl) CreateAuthUser(ctx context.Context, tx *gorm.DB, authUsers *model.AuthUsers) (model.AuthUsers, error) {
	result := tx.WithContext(ctx).Create(&authUsers)
	if result.Error != nil {
		return model.AuthUsers{}, result.Error
	}

	return *authUsers, nil
}

// CreateVerifyEmail implements AuthRepository.
func (a AuthRepositoryImpl) CreateVerifyEmail(ctx context.Context, verifyEmail *model.VerifyEmail) (model.VerifyEmail, error) {
	result := a.db.WithContext(ctx).Create(&verifyEmail)
	if result.Error != nil {
		return model.VerifyEmail{}, result.Error
	}

	return *verifyEmail, nil
}

// FindAuthSessionsById implements AuthRepository.
func (a AuthRepositoryImpl) FindAuthSessionsById(ctx context.Context, id string) (model.AuthSessions, error) {
	session := model.AuthSessions{}
	err := a.db.WithContext(ctx).Model(&model.AuthSessions{}).Take(&session, "id =?", id).Error

	if err != nil {
		return model.AuthSessions{}, err
	}
	return session, nil
}

// FindAuthUserById implements AuthRepository.
func (a AuthRepositoryImpl) FindAuthUserById(ctx context.Context, username string) (model.AuthUsers, error) {
	authUser := model.AuthUsers{}
	err := a.db.WithContext(ctx).Model(&model.AuthUsers{}).Take(&authUser, "username =?", username).Error

	if err != nil {
		return model.AuthUsers{}, err
	}
	return authUser, nil
}

// FindVerifyEmailById implements AuthRepository.
func (a AuthRepositoryImpl) FindVerifyEmailById(ctx context.Context, id int64) (model.VerifyEmail, error) {
	verifyEmail := model.VerifyEmail{}
	err := a.db.WithContext(ctx).Model(&model.VerifyEmail{}).Take(&verifyEmail, "id =?", id).Error

	if err != nil {
		return model.VerifyEmail{}, err
	}
	return verifyEmail, nil
}

// UpdateAuthUser implements AuthRepository.
func (a AuthRepositoryImpl) UpdateAuthUser(ctx context.Context, tx *gorm.DB, authUsers *model.AuthUsers) (model.AuthUsers, error) {
	result := tx.WithContext(ctx).Save(&authUsers)

	if result.Error != nil {
		return model.AuthUsers{}, result.Error
	}

	return *authUsers, nil
}

// UpdateVerifyEmail implements AuthRepository.
func (a AuthRepositoryImpl) UpdateVerifyEmail(ctx context.Context, tx *gorm.DB, verifyEmail *model.VerifyEmail) (model.VerifyEmail, error) {
	result := tx.WithContext(ctx).Save(&verifyEmail)

	if result.Error != nil {
		return model.VerifyEmail{}, result.Error
	}

	return *verifyEmail, nil
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return AuthRepositoryImpl{db: db}
}
