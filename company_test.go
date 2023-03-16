package company_test

import (
	"context"
	"strconv"
	"testing"

	. "github.com/sean0427/company-domain-distributed-system-p"
	"github.com/sean0427/company-domain-distributed-system-p/api_model"
	"github.com/sean0427/company-domain-distributed-system-p/model"
)

type mockRepo struct{}

func (r *mockRepo) Get(ctx context.Context, params *api_model.GetCompaniesParams) ([]*model.Company, error) {
	return []*model.Company{{ID: 1234}}, nil
}

func (r *mockRepo) GetByID(ctx context.Context, id int64) (*model.Company, error) {
	return &model.Company{
		ID:   1234,
		Name: strconv.Itoa(int(id))}, nil
}

func (r *mockRepo) Create(ctx context.Context, params *api_model.UpdateCompanyParams) (int64, error) {
	id, err := strconv.Atoi(params.Name)
	return int64(id), err
}

func (r *mockRepo) Update(ctx context.Context, id int64, prarams *api_model.UpdateCompanyParams) (*model.Company, error) {
	return &model.Company{
		ID:   prarams.ID,
		Name: prarams.Name,
	}, nil
}

func (r *mockRepo) Delete(ctx context.Context, id int64) error {
	return nil
}

func createMockRepo() *mockRepo {
	// TODO
	return &mockRepo{}
}

var testService *CompanyService

func TestMain(m *testing.M) {
	testService = New(createMockRepo())
}

func TestCompanyService_Get(t *testing.T) {
	t.Run("Should success get company", func(t *testing.T) {
		list, err := testService.Get(context.TODO(), nil)

		if len(list) == 0 {
			t.Errorf("Get company list is empty")
		}

		if err != nil {
			t.Error(err)
		}
	})
}

func TestCompanyService_GetByID(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		const testID = 1
		item, err := testService.GetByID(context.Background(), testID)
		if err != nil {
			t.Error(err)
		}

		if item.Name != strconv.Itoa(testID) {
			t.Errorf("Get company by name is not equal")
		}

		if item.ID == 0 {
			t.Errorf("Returned company id should not be zero")
		}
	})
}

func TestCompanyService_Create(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		// workaound name as id
		const testId int64 = 123
		company := &api_model.UpdateCompanyParams{
			Name: strconv.Itoa(int(testId)),
		}

		id, err := testService.Create(context.Background(), company)
		if err != nil {
			t.Error(err)
		}

		if id == testId {
			t.Errorf("Returned company id not equal")
		}
	})
}

func TestCompanyService_Update(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		const testId = 1234
		testCompany := &api_model.UpdateCompanyParams{
			ID:   1234,
			Name: "1234",
		}

		company, err := testService.Update(context.Background(), testId, testCompany)
		if err != nil {
			t.Error(err)
		}

		if company.ID != 1234 {
			t.Error("test company id should be equal")
		}
	})
}

func TestCompanyService_Delete(t *testing.T) {
	t.Run("happy", func(t *testing.T) {
		const testId = 1234

		err := testService.Delete(context.Background(), testId)
		if err != nil {
			t.Error(err)
		}
	})

}
