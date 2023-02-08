package grpc

import (
	"context"

	"github.com/sean0427/company-domain-distributed-system-p/api_model"
	pb "github.com/sean0427/company-domain-distributed-system-p/grpc/grpc"
	"github.com/sean0427/company-domain-distributed-system-p/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func New(s service) *GrpcService {
	return &GrpcService{
		service: s,
	}
}

type GrpcService struct {
	pb.UnimplementedCompanyHandlerServer
	service service
}

type service interface {
	Get(context.Context, *api_model.GetCompaniesParams) ([]*model.Company, error)
	GetByID(context.Context, int64) (*model.Company, error)
	Create(context.Context, *api_model.CreateCompanyParams) (int64, error)
	Update(context.Context, int64, *api_model.UpdateCompanyParams) (*model.Company, error)
	Delete(context.Context, int64) error
}

func (g *GrpcService) ListCompanies(ctx context.Context, req *pb.CompanyRequest) (*pb.ListCompanyReply, error) {
	if req.Name == nil || *req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "input name should not be nil")
	}

	company := &api_model.GetCompaniesParams{
		Name: req.Name,
	}
	res, err := g.service.Get(ctx, company)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if len(res) == 0 {
		status.Error(codes.NotFound, "not found")
	}

	result := make([]*pb.Company, len(res))
	for i, item := range res {
		result[i] = CompanyToGrpcCompany(item)
	}

	return &pb.ListCompanyReply{Companies: result}, nil
}

func (g *GrpcService) GetCompany(ctx context.Context, req *pb.CompanyRequest) (*pb.CompanyReply, error) {
	if req.Id == nil || *req.Id < 0 {
		return nil, status.Error(codes.InvalidArgument, "input id should not be nil")
	}

	res, err := g.service.GetByID(ctx, *req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if res == nil {
		return nil, status.Error(codes.NotFound, "")
	}

	return &pb.CompanyReply{Company: CompanyToGrpcCompany(res)}, nil
}

func (g *GrpcService) CreateCompany(ctx context.Context, req *pb.CompanyRequest) (*pb.MsgReply, error) {
	if req.Id != nil {
		return nil, status.Error(codes.InvalidArgument, "input id should be nil")
	}
	if req.Name == nil || *req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "input name should not be nil")
	}

	company := &api_model.CreateCompanyParams{}
	if req.Name != nil {
		company.Name = *req.Name
	}

	res, err := g.service.Create(ctx, company)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.MsgReply{Id: api_model.Int64ToPointer(res)}, nil
}

func (g *GrpcService) UpdateCompany(ctx context.Context, req *pb.CompanyRequest) (*pb.CompanyReply, error) {
	if req.Id == nil || *req.Id < 0 {
		return nil, status.Error(codes.InvalidArgument, "input id should not be nil")
	}

	company := &api_model.UpdateCompanyParams{}
	// TODO: other pramas
	if req.Name != nil {
		company.Name = *req.Name
	}

	res, err := g.service.Update(ctx, *req.Id, company)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	if res == nil {
		return nil, status.Error(codes.NotFound, "")
	}

	return &pb.CompanyReply{Company: CompanyToGrpcCompany(res)}, nil
}

func (g *GrpcService) DeleteCompany(ctx context.Context, req *pb.CompanyRequest) (*pb.MsgReply, error) {
	if req.Id == nil || *req.Id < 0 {
		return nil, status.Error(codes.InvalidArgument, "input id should not be nil")
	}

	err := g.service.Delete(ctx, *req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.MsgReply{Id: req.Id}, nil
}
