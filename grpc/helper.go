package grpc

import (
	"github.com/sean0427/company-domain-distributed-system-p/grpc/grpc"
	"github.com/sean0427/company-domain-distributed-system-p/model"
)

func CompanyToGrpcCompany(item *model.Company) *grpc.Company {
	if item == nil {
		return nil
	}
	return &grpc.Company{
		Id:    item.ID,
		Name:  item.Name,
		Email: item.Email,
	}
}
