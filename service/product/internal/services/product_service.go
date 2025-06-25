package services

import (
	"context"
	"e-commerce-microservice/product/internal/config"
	db "e-commerce-microservice/product/internal/db/sqlc"
	"e-commerce-microservice/product/internal/pb"
	"e-commerce-microservice/product/internal/token"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
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
	var response *pb.ProductDetail

	err := p.Store.ExecTx(ctx, func(q *db.Queries) error {
		Category, err := q.GetCategories(ctx, req.GetCategoryId())
		if err != nil {
			return fmt.Errorf("Failed get Category")
		}

		productParam := db.CreateProductParams{
			CategoryID:  Category.ID,
			Name:        req.Name,
			Description: req.Description,
			Images:      []string{},
			Rating:      0,
			Price:       req.Price,
		}
		productResult, err := q.CreateProduct(ctx, productParam)

		if err != nil {
			return fmt.Errorf("Failed to create product")
		}

		var resultColorVarians []*pb.ColorVarian

		for i := range req.ColorVarians {
			colorVarianParam := db.CreateColorVarianProductParams{
				ProductID: productResult.ID,
				Name:      req.ColorVarians[i].Name,
				Color:     req.ColorVarians[i].Color,
				Images:    []string{},
			}

			colorVarianResult, err := q.CreateColorVarianProduct(ctx, colorVarianParam)
			if err != nil {
				return fmt.Errorf("Failed to create color varian ke %d", i)
			}

			var resultSizeVarians []*pb.SizeVarian
			sizeVarians := req.ColorVarians[i].SizeVarians
			for j := range sizeVarians {
				sizeParam := db.CreateSizeVarianProductParams{
					ColorVarianID: colorVarianResult.ID,
					Size:          sizeVarians[j].Size,
					Stock:         sizeVarians[j].Stock,
				}
				size, err := q.CreateSizeVarianProduct(ctx, sizeParam)
				if err != nil {
					return fmt.Errorf("Failed to create size varian ke %d", j)
				}
				resultSize := pb.SizeVarian{
					Id:            size.ID,
					ColorVarianId: size.ColorVarianID,
					Size:          size.Size,
					Stock:         size.Stock,
					CreatedAt:     timestamppb.New(size.CreatedAt),
					UpdatedAt:     timestamppb.New(size.UpdatedAt),
				}
				resultSizeVarians = append(resultSizeVarians, &resultSize)
			}
			resultColor := &pb.ColorVarian{
				Id:          colorVarianResult.ID,
				ProductId:   colorVarianResult.ProductID,
				Name:        colorVarianResult.Name,
				Color:       colorVarianResult.Color,
				Images:      colorVarianResult.Images,
				SizeVarians: resultSizeVarians,
				CreatedAt:   timestamppb.New(colorVarianResult.CreatedAt),
				UpdatedAt:   timestamppb.New(colorVarianResult.UpdatedAt),
			}
			resultColorVarians = append(resultColorVarians, resultColor)
		}

		response = &pb.ProductDetail{
			Id:           productResult.ID,
			Name:         productResult.Name,
			Price:        productResult.Price,
			Images:       productResult.Images,
			Rating:       productResult.Rating,
			ColorVarians: resultColorVarians,
			Category: &pb.Category{
				Id:        productResult.CategoryID,
				Name:      Category.Name,
				Icon:      Category.Icon,
				UpdatedAt: timestamppb.New(Category.UpdatedAt),
				CreatedAt: timestamppb.New(Category.CreatedAt),
			},
			CreatedAt: timestamppb.New(productResult.CreatedAt),
			UpdatedAt: timestamppb.New(productResult.UpdatedAt),
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
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
