package utils

import (
	"e-commerce-microservice/user/internal/model"
	"e-commerce-microservice/user/internal/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserAddress(adress model.UserAddresses) *pb.UserAddress {
	return &pb.UserAddress{
		Id:             adress.ID,
		Username:       adress.Username,
		Label:          adress.Label,
		RecipientName:  adress.RecipientName,
		RecipientPhone: adress.RecipientPhone,
		AddressLine:    adress.AddressLine,
		City:           adress.City,
		Province:       adress.Province,
		PostalCode:     adress.PostalCode,
		IsSelected:     adress.IsSelected,
		UpdatedAt:      timestamppb.New(adress.UpdatedAt),
		CreatedAt:      timestamppb.New(adress.CreatedAt),
	}
}
