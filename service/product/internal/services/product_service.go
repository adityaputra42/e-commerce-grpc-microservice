package services

import (
	"context"
	"e-commerce-microservice/product/internal/config"
	db "e-commerce-microservice/product/internal/db/sqlc"
	"e-commerce-microservice/product/internal/pb"
	"e-commerce-microservice/product/internal/token"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductDetail, error)
	GetProduct(ctx context.Context, req *pb.GetByIDRequest) (*pb.ProductDetail, error)
	GetListProduct(ctx context.Context, req *pb.PaginationRequest) (*pb.ProductList, error)
	UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductDetail, error)
	DeleteProduct(ctx context.Context, req *pb.GetByIDRequest) error
}

type ProductServiceImpl struct {
	Config     config.Configuration
	Store      db.Store
	TokenMaker token.Maker
}

// CreateProduct implements ProductService.
func (p ProductServiceImpl) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductDetail, error) {
	panic("unimplemented")
}

// DeleteProduct implements ProductService.
func (p ProductServiceImpl) DeleteProduct(ctx context.Context, req *pb.GetByIDRequest) error {
	panic("unimplemented")
}

// GetListProduct implements ProductService.
func (p ProductServiceImpl) GetListProduct(ctx context.Context, req *pb.PaginationRequest) (*pb.ProductList, error) {
	panic("unimplemented")
}

// GetProduct implements ProductService.
func (p ProductServiceImpl) GetProduct(ctx context.Context, req *pb.GetByIDRequest) (*pb.ProductDetail, error) {
	panic("unimplemented")
}

// UpdateProduct implements ProductService.
func (p ProductServiceImpl) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductDetail, error) {
	panic("unimplemented")
}

func NewProductService(Config config.Configuration,
	Store db.Store,
	TokenMaker token.Maker) ProductService {
	return ProductServiceImpl{Config: Config, Store: Store, TokenMaker: TokenMaker}
}
