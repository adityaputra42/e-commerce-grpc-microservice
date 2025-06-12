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
	value, err := c.Store.GetCategories(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get category")
	}
	resp := &pb.Category{Id: value.ID, Name: value.Name, Icon: value.Icon, UpdatedAt: timestamppb.New(value.UpdatedAt), CreatedAt: timestamppb.New(value.CreatedAt)}
	return resp, nil
}

// GetListCategory implements CategoryService.
func (c CategoryServiceImpl) GetListCategory(ctx context.Context) (*pb.CategoryList, error) {

	listCategory, err := c.Store.ListCategories(ctx, db.ListCategoriesParams{})

	categories := []*pb.Category{}
	if err != nil {
		return nil, fmt.Errorf("Failed to get list category")
	}

	for _, value := range listCategory {

		categories = append(categories, &pb.Category{Id: value.ID, Name: value.Name, Icon: value.Icon, UpdatedAt: timestamppb.New(value.UpdatedAt), CreatedAt: timestamppb.New(value.CreatedAt)})

	}

	return &pb.CategoryList{
		Categories: categories,
	}, nil

}

// UpdateCategory implements CategoryService.
func (c CategoryServiceImpl) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.Category, error) {
	_, err := utils.AuthorizationUser(ctx, []string{utils.AdminRole}, c.TokenMaker)

	if err != nil {
		return nil, utils.UnauthenticatedError(err)
	}

	category, err := c.Store.GetCategoriesForUpdate(ctx, req.GetId())
	if err != nil {
		return nil, fmt.Errorf("Failed to get category")
	}

	param := db.UpdateCategoriesParams{
		ID:   category.ID,
		Name: req.Name,
		Icon: req.Icon,
	}

	newCategory, err := c.Store.UpdateCategories(ctx, param)

	if err != nil {
		return nil, fmt.Errorf("Failed to update category")
	}
	resp := &pb.Category{Id: newCategory.ID, Name: newCategory.Name, Icon: newCategory.Icon, UpdatedAt: timestamppb.New(newCategory.UpdatedAt), CreatedAt: timestamppb.New(newCategory.CreatedAt)}
	return resp, nil

}

func NewCategoryService(Config config.Configuration,
	Store db.Store,
	TokenMaker token.Maker) CategoryService {
	return CategoryServiceImpl{Config: Config, Store: Store, TokenMaker: TokenMaker}
}
