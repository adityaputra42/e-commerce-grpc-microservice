package services

import (
	"context"
	"e-commerce-microservice/product/internal/config"
	db "e-commerce-microservice/product/internal/db/sqlc"
	"e-commerce-microservice/product/internal/pb"
	"e-commerce-microservice/product/internal/token"
)

type ColorVarianService interface {
	AddColorVarian(ctx context.Context, req *pb.CreateColorVarianRequest) (*pb.ColorVarian, error)
	GetColorVarian(ctx context.Context, req *pb.GetByIDRequest) (*pb.ColorVarian, error)
	GetListColorVarian(ctx context.Context) (*pb.ColorVarianList, error)
	UpdateColorVarian(ctx context.Context, req *pb.UpdateColorVarianRequest) (*pb.ColorVarian, error)
	DeleteColorVarian(ctx context.Context, req *pb.GetByIDRequest) error
}

type ColorVarianServiceImpl struct {
	Config     config.Configuration
	Store      db.Store
	TokenMaker token.Maker
}

// AddColorVarian implements ColorVarianService.
func (c ColorVarianServiceImpl) AddColorVarian(ctx context.Context, req *pb.CreateColorVarianRequest) (*pb.ColorVarian, error) {
	panic("unimplemented")
}

// DeleteColorVarian implements ColorVarianService.
func (c ColorVarianServiceImpl) DeleteColorVarian(ctx context.Context, req *pb.GetByIDRequest) error {
	panic("unimplemented")
}

// GetColorVarian implements ColorVarianService.
func (c ColorVarianServiceImpl) GetColorVarian(ctx context.Context, req *pb.GetByIDRequest) (*pb.ColorVarian, error) {
	panic("unimplemented")
}

// GetListColorVarian implements ColorVarianService.
func (c ColorVarianServiceImpl) GetListColorVarian(ctx context.Context) (*pb.ColorVarianList, error) {
	panic("unimplemented")
}

// UpdateColorVarian implements ColorVarianService.
func (c ColorVarianServiceImpl) UpdateColorVarian(ctx context.Context, req *pb.UpdateColorVarianRequest) (*pb.ColorVarian, error) {
	panic("unimplemented")
}

func NewColorVarianService(Config config.Configuration,
	Store db.Store,
	TokenMaker token.Maker) ColorVarianService {
	return ColorVarianServiceImpl{Config: Config, Store: Store, TokenMaker: TokenMaker}
}
