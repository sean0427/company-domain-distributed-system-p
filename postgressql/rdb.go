package postgressql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/sean0427/company-domain-distributed-system-p/api_model"
	"github.com/sean0427/company-domain-distributed-system-p/model"
)

const topicName = "company"

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *repository {
	return &repository{
		db: db,
	}
}

func getUserFromContext(ctx context.Context) (string, error) {
	user := ctx.Value("user")
	if user == nil {
		return "default", nil // TODO
	}

	return user.(string), nil
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
	user, _ := getUserFromContext(ctx)
	company := model.Company{
		Name:      params.Name,
		Email:     params.Email,
		Address:   params.Address,
		Contact:   params.Contact,
		UpdatedBy: user,
		CreatedBy: user,
	}

	err := TransactionWithOutboxMsg(ctx, r.db, &company, topicName, func(tx *gorm.DB) (int64, error) {
		result := tx.Model(&company).Create(&company)
		return company.ID, result.Error
	})

	if err != nil {
		return 0, err
	}

	return company.ID, nil
}

func (r *repository) Update(ctx context.Context, id int64, params *api_model.UpdateCompanyParams) (*model.Company, error) {
	user, _ := getUserFromContext(ctx)
	company := model.Company{
		ID:        id,
		Name:      params.Name,
		Email:     params.Email,
		Address:   params.Address,
		Contact:   params.Contact,
		UpdatedBy: user,
	}
	err := TransactionWithOutboxMsg(ctx, r.db, &company, topicName, func(tx *gorm.DB) (int64, error) {
		result := tx.Model(&company).Where("id = ?", params.ID).Save(&company)
		if result.Error != nil {
			return 0, result.Error
		}
		if result.RowsAffected == 0 {
			return 0, errors.New("not found")
		}
		return params.ID, nil
	})
	if err != nil {
		return nil, err
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
