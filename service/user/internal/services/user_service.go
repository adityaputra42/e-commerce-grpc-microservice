package services

import (
	"context"
	"e-commerce-microservice/user/internal/model"
	"e-commerce-microservice/user/internal/pb"
	"e-commerce-microservice/user/internal/repository"
	"e-commerce-microservice/user/internal/token"
	"e-commerce-microservice/user/internal/utils"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type UserService interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error)
	FindUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error)
	CreateUserAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.UserAddressResponse, error)
	FindUserAddresses(ctx context.Context, req *pb.GetUserAddressesRequest) (*pb.UserAddressesResponse, error)
	UpdateUserAddresses(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.UserAddressResponse, error)
}
type UserServiceImpl struct {
	tokenMaker token.Maker
	db         *gorm.DB
	repo       repository.UserRepository
}

// CreateUser implements UserService.
func (u *UserServiceImpl) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	var user pb.UserResponse

	err := u.db.Transaction(func(tx *gorm.DB) error {
		authPayload, err := utils.AuthorizationUser(ctx, u.tokenMaker)

		if err != nil {
			return utils.UnauthenticatedError(err)
		}
		if authPayload.Username != req.Username {
			return fmt.Errorf("username mismatch")
		}
		userReq := model.Users{
			Username:    authPayload.Username,
			FullName:    req.GetFullName(),
			PhoneNumber: req.GetPhoneNumber(),
		}
		userResp, err := u.repo.CreateUser(ctx, tx, &userReq)
		if err != nil {
			return fmt.Errorf("Failed to create user")
		}

		user = pb.UserResponse{
			User: &pb.User{
				Username:    userResp.Username,
				FullName:    userResp.FullName,
				PhoneNumber: userResp.PhoneNumber,
				UpdatedAt:   timestamppb.New(userResp.UpdatedAt),
				CreatedAt:   timestamppb.New(userResp.CreatedAt),
			},
		}
		return nil
	})

	if err != nil {
		return &pb.UserResponse{}, err
	}

	return &user, nil
}

// CreateUserAddress implements UserService.
func (u *UserServiceImpl) CreateUserAddress(ctx context.Context, req *pb.CreateAddressRequest) (*pb.UserAddressResponse, error) {
	var user pb.UserAddressResponse

	err := u.db.Transaction(func(tx *gorm.DB) error {
		authPayload, err := utils.AuthorizationUser(ctx, u.tokenMaker)

		if err != nil {
			return utils.UnauthenticatedError(err)
		}
		if authPayload.Username != req.Username {
			return fmt.Errorf("username mismatch")
		}
		addressReq := model.UserAddresses{
			Username:       authPayload.Username,
			Label:          req.GetLabel(),
			RecipientName:  req.GetRecipientName(),
			RecipientPhone: req.GetRecipientPhone(),
			AddressLine:    req.GetAddressLine(),
			City:           req.GetCity(),
			Province:       req.GetProvince(),
			PostalCode:     req.GetPostalCode(),
			IsSelected:     true,
		}
		addressRes, err := u.repo.CreateUserAddress(ctx, tx, &addressReq)
		if err != nil {
			return fmt.Errorf("Failed to create user Address")
		}

		user = pb.UserAddressResponse{
			Address: utils.ToUserAddress(addressRes),
		}
		return nil
	})

	if err != nil {
		return &pb.UserAddressResponse{}, err
	}

	return &user, nil

}

// FindUser implements UserService.
func (u *UserServiceImpl) FindUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	var user pb.UserResponse

	authPayload, err := utils.AuthorizationUser(ctx, u.tokenMaker)

	if err != nil {
		return &pb.UserResponse{}, utils.UnauthenticatedError(err)
	}
	if authPayload.Username != req.Username {
		return &pb.UserResponse{}, fmt.Errorf("username mismatch")
	}

	res, err := u.repo.FindUserById(ctx, req.GetUsername())

	if err != nil {
		return &pb.UserResponse{}, fmt.Errorf("failed to get user")
	}

	user = pb.UserResponse{
		User: &pb.User{
			Username:    res.Username,
			FullName:    res.FullName,
			PhoneNumber: res.PhoneNumber,
			UpdatedAt:   timestamppb.New(res.UpdatedAt),
			CreatedAt:   timestamppb.New(res.CreatedAt),
		},
	}

	return &user, nil
}

// FindUserAddresses implements UserService.
func (u *UserServiceImpl) FindUserAddresses(ctx context.Context, req *pb.GetUserAddressesRequest) (*pb.UserAddressesResponse, error) {

	var userAddress pb.UserAddressesResponse

	authPayload, err := utils.AuthorizationUser(ctx, u.tokenMaker)

	if err != nil {
		return &pb.UserAddressesResponse{}, utils.UnauthenticatedError(err)
	}
	if authPayload.Username != req.Username {
		return &pb.UserAddressesResponse{}, fmt.Errorf("username mismatch")
	}

	addresses, err := u.repo.FindUserAdrressesByUsername(ctx, req.Username)

	if err != nil {
		return &pb.UserAddressesResponse{}, fmt.Errorf("failed to get user addresses")
	}

	newAddress := []*pb.UserAddress{}

	for _, v := range addresses {
		newAddress = append(newAddress, utils.ToUserAddress(v))
	}

	userAddress = pb.UserAddressesResponse{
		Addresses: newAddress,
	}

	return &userAddress, nil
}

// UpdateUserAddresses implements UserService.
func (u *UserServiceImpl) UpdateUserAddresses(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.UserAddressResponse, error) {
	var updatedAddress pb.UserAddressResponse

	err := u.db.Transaction(func(tx *gorm.DB) error {

		_, err := utils.AuthorizationUser(ctx, u.tokenMaker)

		if err != nil {
			return utils.UnauthenticatedError(err)
		}

		address, err := u.repo.FindUserAdrressById(ctx, req.Id)

		if err != nil {
			return fmt.Errorf("failed to get user address")
		}

		request := model.UserAddresses{
			ID:             address.ID,
			Username:       address.Username,
			Label:          req.Label,
			RecipientName:  req.RecipientName,
			RecipientPhone: req.RecipientPhone,
			AddressLine:    req.AddressLine,
			City:           req.City,
			Province:       req.Province,
			PostalCode:     req.PostalCode,
			IsSelected:     req.IsSelected,
		}

		addressNew, err := u.repo.UpdateUserAddress(ctx, tx, &request)

		if err != nil {
			return fmt.Errorf("failed to update user address")
		}
		updatedAddress = pb.UserAddressResponse{
			Address: utils.ToUserAddress(addressNew),
		}
		return nil
	})

	if err != nil {
		return &pb.UserAddressResponse{}, err
	}

	return &updatedAddress, nil
}

func NewUserService(tokenMaker token.Maker,
	db *gorm.DB,
	repo repository.UserRepository) UserService {
	return &UserServiceImpl{tokenMaker: tokenMaker, db: db, repo: repo}
}
