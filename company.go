package company

import (
	"context"

	"github.com/sean0427/company-domain-distributed-system-p/api_model"
	"github.com/sean0427/company-domain-distributed-system-p/model"
)

type repository interface {
	Get(ctx context.Context, params *api_model.GetCompaniesParams) ([]*model.Company, error)
	GetByID(ctx context.Context, id int64) (*model.Company, error)
	Create(ctx context.Context, company *api_model.CreateCompanyParams) (int64, error)
	Update(ctx context.Context, id int64, company *api_model.UpdateCompanyParams) (*model.Company, error)
	Delete(ctx context.Context, id int64) error
}

type CompanyService struct {
	repo repository
}

func New(repo repository) *CompanyService {
	return &CompanyService{
		repo: repo,
	}
}

func (s *CompanyService) Get(ctx context.Context, params *api_model.GetCompaniesParams) ([]*model.Company, error) {
	return s.repo.Get(ctx, params)
}

func (s *CompanyService) GetByID(ctx context.Context, id int64) (*model.Company, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CompanyService) Create(ctx context.Context, params *api_model.CreateCompanyParams) (int64, error) {
	return s.repo.Create(ctx, params)
}

func (s *CompanyService) Update(ctx context.Context, id int64, params *api_model.UpdateCompanyParams) (*model.Company, error) {
	return s.repo.Update(ctx, id, params)
}

func (s *CompanyService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
