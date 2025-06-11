package services

import (
	"context"
	"e-commerce-microservice/product/internal/config"
	db "e-commerce-microservice/product/internal/db/sqlc"
	"e-commerce-microservice/product/internal/pb"
	"e-commerce-microservice/product/internal/token"
)

type SizeVarianService interface {
	AddSizeVarian(ctx context.Context, req *pb.CreateSizeVarianRequest) (*pb.SizeVarian, error)
	GetSizeVarian(ctx context.Context, req *pb.GetByIDRequest) (*pb.SizeVarian, error)
	GetListSizeVarian(ctx context.Context) (*pb.SizeVarianList, error)
	UpdateSizeVarian(ctx context.Context, req *pb.UpdateSizeVarianRequest) (*pb.SizeVarian, error)
	DeleteSizeVarian(ctx context.Context, req *pb.GetByIDRequest) error
}

type SizeVarianServiceImpl struct {
	Config     config.Configuration
	Store      db.Store
	TokenMaker token.Maker
}

// AddSizeVarian implements SizeVarianService.
func (s SizeVarianServiceImpl) AddSizeVarian(ctx context.Context, req *pb.CreateSizeVarianRequest) (*pb.SizeVarian, error) {
	panic("unimplemented")
}

// DeleteSizeVarian implements SizeVarianService.
func (s SizeVarianServiceImpl) DeleteSizeVarian(ctx context.Context, req *pb.GetByIDRequest) error {
	panic("unimplemented")
}

// GetListSizeVarian implements SizeVarianService.
func (s SizeVarianServiceImpl) GetListSizeVarian(ctx context.Context) (*pb.SizeVarianList, error) {
	panic("unimplemented")
}

// GetSizeVarian implements SizeVarianService.
func (s SizeVarianServiceImpl) GetSizeVarian(ctx context.Context, req *pb.GetByIDRequest) (*pb.SizeVarian, error) {
	panic("unimplemented")
}

// UpdateSizeVarian implements SizeVarianService.
func (s SizeVarianServiceImpl) UpdateSizeVarian(ctx context.Context, req *pb.UpdateSizeVarianRequest) (*pb.SizeVarian, error) {
	panic("unimplemented")
}

func NewSizeVarianService(Config config.Configuration,
	Store db.Store,
	TokenMaker token.Maker) SizeVarianService {
	return SizeVarianServiceImpl{Config: Config, Store: Store, TokenMaker: TokenMaker}
}
