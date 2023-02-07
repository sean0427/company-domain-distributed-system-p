package postgressql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/sean0427/company-domain-distributed-system-p/api_model"
	"github.com/sean0427/company-domain-distributed-system-p/model"
)

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Get(ctx context.Context, params *api_model.GetCompaniesParams) ([]*model.Company, error) {
	var companies []*model.Company
	tx := r.db.WithContext(ctx)

	if params != nil && (*params).Name != nil {
		tx = tx.Where("name = ?", params.Name)
	}

	result := tx.Find(&companies)

	return companies, result.Error
}

func (r *repository) GetByID(ctx context.Context, id int64) (*model.Company, error) {
	var company model.Company
	tx := r.db.WithContext(ctx)

	result := tx.Where("id = ?", id).Find(&company)
	return &company, result.Error
}

func (r *repository) Create(ctx context.Context, params *api_model.CreateCompanyParams) (int64, error) {
	company := model.Company{
		Name:     params.Name,
		Email:    params.Email,
		Password: params.Password,
	}

	tx := r.db.WithContext(ctx)
	result := tx.Model(&company).Create(&company)
	if result.Error != nil {
		return 0, result.Error
	}

	return company.ID, nil
}

func (r *repository) Update(ctx context.Context, id int64, params *api_model.UpdateCompanyParams) (*model.Company, error) {
	company := model.Company{
		ID:       id,
		Name:     params.Name,
		Email:    params.Email,
		Password: params.Password,
	}
	tx := r.db.WithContext(ctx)

	result := tx.Model(&company).Where("id = ?", params.ID).Save(&company)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("not found")
	}

	return &company, nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	tx := r.db.WithContext(ctx)

	result := tx.Delete(&model.Company{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("not found")
	}

	return nil
}

func (r *repository) ExamCompanyPassword(ctx context.Context, name, password string) (bool, error) {
	tx := r.db.WithContext(ctx)

	var count int64
	result := tx.Model(&model.Company{}).
		Where("name =? and password=?", name, password).
		Count(&count)

	if count == 1 {
		return true, nil
	}

	return false, result.Error
}
