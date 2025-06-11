package services

import (
	"context"
	"e-commerce-microservice/product/internal/config"
	db "e-commerce-microservice/product/internal/db/sqlc"
	"e-commerce-microservice/product/internal/pb"
	"e-commerce-microservice/product/internal/token"
	"e-commerce-microservice/product/internal/utils"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.Category, error)
	GetCategory(ctx context.Context, req *pb.GetByIDRequest) (*pb.Category, error)
	GetListCategory(ctx context.Context) (*pb.CategoryList, error)
	UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.Category, error)
	DeleteCategory(ctx context.Context, req *pb.GetByIDRequest) error
}

type CategoryServiceImpl struct {
	Config     config.Configuration
	Store      db.Store
	TokenMaker token.Maker
}

// CreateCategory implements CategoryService.
func (c CategoryServiceImpl) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.Category, error) {

	_, err := utils.AuthorizationUser(ctx, []string{utils.AdminRole}, c.TokenMaker)

	if err != nil {
		return nil, utils.UnauthenticatedError(err)
	}
	param := db.CreateCategoriesParams{
		Name: req.Name,
		Icon: req.Icon,
	}

	Category, err := c.Store.CreateCategories(ctx, param)
	if err != nil {
		return nil, err
	}
	return &pb.Category{
		Id:        Category.ID,
		Name:      Category.Name,
		Icon:      Category.Icon,
		UpdatedAt: timestamppb.New(Category.UpdatedAt),
		CreatedAt: timestamppb.New(Category.CreatedAt),
	}, nil
}

// DeleteCategory implements CategoryService.
func (c CategoryServiceImpl) DeleteCategory(ctx context.Context, req *pb.GetByIDRequest) error {
	_, err := utils.AuthorizationUser(ctx, []string{utils.AdminRole}, c.TokenMaker)

	if err != nil {
		return utils.UnauthenticatedError(err)
	}
	err = c.Store.DeleteCategories(ctx, req.Id)
	if err != nil {
		return fmt.Errorf("Failed to delete category")
	}
	return nil
}

// GetCategory implements CategoryService.
func (c CategoryServiceImpl) GetCategory(ctx context.Context, req *pb.GetByIDRequest) (*pb.Category, error) {
	panic("unimplemented")
}

// GetListCategory implements CategoryService.
func (c CategoryServiceImpl) GetListCategory(ctx context.Context) (*pb.CategoryList, error) {
	panic("unimplemented")
}

// UpdateCategory implements CategoryService.
func (c CategoryServiceImpl) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.Category, error) {
	panic("unimplemented")
}

func NewCategoryService(Config config.Configuration,
	Store db.Store,
	TokenMaker token.Maker) CategoryService {
	return CategoryServiceImpl{Config: Config, Store: Store, TokenMaker: TokenMaker}
}
