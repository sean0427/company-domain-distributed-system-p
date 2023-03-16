package grpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sean0427/company-domain-distributed-system-p/api_model"
	. "github.com/sean0427/company-domain-distributed-system-p/grpc"
	pb "github.com/sean0427/company-domain-distributed-system-p/grpc/grpc"
	mock "github.com/sean0427/company-domain-distributed-system-p/mock"
	"github.com/sean0427/company-domain-distributed-system-p/model"
)

var testListCompanyCasese = []struct {
	name              string
	request           *pb.CompanyQuery
	expectTime        int
	returnedCompanies []*model.Company
	returnedErr       error
	wantErr           bool
}{
	{
		name: "happy",
		request: &pb.CompanyQuery{
			Name: api_model.StringToPointer("any"),
		},
		expectTime: 1,
		returnedCompanies: []*model.Company{
			{
				ID:   1,
				Name: "any",
			},
			{
				ID:   2,
				Name: "any2",
			},
		},
		returnedErr: nil,
		wantErr:     false,
	},
	{
		name: "errpr - not found",
		request: &pb.CompanyQuery{
			Name: api_model.StringToPointer("any"),
		},
		expectTime:        1,
		returnedCompanies: []*model.Company{},
		returnedErr:       nil,
		wantErr:           true,
	}, {
		name:              "error - no name",
		request:           &pb.CompanyQuery{},
		expectTime:        0,
		returnedCompanies: nil,
		returnedErr:       nil,
		wantErr:           true,
	},
	{
		name: "error - servce return err",
		request: &pb.CompanyQuery{
			Name: api_model.StringToPointer("any"),
		},
		expectTime:        1,
		returnedCompanies: nil,
		returnedErr:       errors.New("any"),
		wantErr:           true,
	},
}

func TestGrpcService_ListCompanies(t *testing.T) {
	for _, c := range testListCompanyCasese {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			m := mock.NewMockservice(ctrl)
			m.EXPECT().
				Get(gomock.Any(), gomock.Eq(&api_model.GetCompaniesParams{Name: c.request.Name})).
				Return(c.returnedCompanies, c.returnedErr).
				Times(c.expectTime)

			grpc := New(m)

			res, err := grpc.ListCompanies(context.Background(), c.request)

			if c.wantErr && err != nil {
				if c.returnedErr != nil && err == nil {
					t.Fatalf("expect err %s", err.Error())
				}
				return
			}

			if v, e := len(res.Companies), len(c.returnedCompanies); v != e {
				t.Fatalf("ecpect %d but %d", v, e)
			}
		})
	}
}

var testGetCompanyCase = []struct {
	name            string
	request         *pb.CompanyRequest
	expectTime      int
	returnedCompany *model.Company
	returnedErr     error
	wantErr         bool
}{
	{
		name: "happy",
		request: &pb.CompanyRequest{
			Id: api_model.Int64ToPointer(11),
		},
		expectTime: 1,
		returnedCompany: &model.Company{
			ID:   1,
			Name: "any",
		},
		returnedErr: nil,
		wantErr:     false,
	},
	{
		name: "error - not found",
		request: &pb.CompanyRequest{
			Id: api_model.Int64ToPointer(2211),
		},
		expectTime:      1,
		returnedCompany: nil,
		returnedErr:     nil,
		wantErr:         true,
	},
	{
		name:            "error - no id",
		request:         &pb.CompanyRequest{},
		expectTime:      0,
		returnedCompany: nil,
		returnedErr:     nil,
		wantErr:         true,
	},
	{
		name: "error - servce return err",
		request: &pb.CompanyRequest{
			Id: api_model.Int64ToPointer(11),
		},
		expectTime:      1,
		returnedCompany: nil,
		returnedErr:     errors.New("any"),
		wantErr:         true,
	},
}

func TestGrpcService_GetCompany(t *testing.T) {
	for _, c := range testGetCompanyCase {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			idMatcher := gomock.Any()
			if c.request != nil && c.request.Id != nil {
				idMatcher = gomock.Eq(*c.request.Id)
			}

			m := mock.NewMockservice(ctrl)
			m.EXPECT().
				GetByID(gomock.Any(), idMatcher).
				Return(c.returnedCompany, c.returnedErr).
				Times(c.expectTime)

			grpc := New(m)

			res, err := grpc.GetCompany(context.Background(), c.request)

			if c.wantErr && err != nil {
				if c.returnedErr != nil && err == nil {
					t.Fatalf("expect err %s", err.Error())
				}
				return
			}

			if res == nil {
				t.Fatal("expect not nil")
			}

			if v, e := res.Company.Id, c.returnedCompany.ID; v != e {
				t.Fatalf("ecpect id: %d but %d", v, e)
			}
		})
	}
}

var testDeleteCompanyCase = []struct {
	name        string
	request     *pb.CompanyRequest
	expectTime  int
	returnedErr error
	wantErr     bool
}{
	{
		name: "happy",
		request: &pb.CompanyRequest{
			Id: api_model.Int64ToPointer(11),
		},
		expectTime:  1,
		returnedErr: nil,
		wantErr:     false,
	},
	{
		name: "error - not found",
		request: &pb.CompanyRequest{
			Id: api_model.Int64ToPointer(2211),
		},
		expectTime:  1,
		returnedErr: nil,
		wantErr:     true,
	}, {
		name:        "error - no id",
		request:     &pb.CompanyRequest{},
		expectTime:  0,
		returnedErr: nil,
		wantErr:     true,
	},
	{
		name: "error - servce return err",
		request: &pb.CompanyRequest{
			Id: api_model.Int64ToPointer(11),
		},
		expectTime:  1,
		returnedErr: errors.New("any"),
		wantErr:     true,
	},
}

func TestGrpcService_DeleteCompany(t *testing.T) {
	for _, c := range testDeleteCompanyCase {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			idMatcher := gomock.Any()
			if c.request != nil && c.request.Id != nil {
				idMatcher = gomock.Eq(*c.request.Id)
			}

			m := mock.NewMockservice(ctrl)
			m.EXPECT().
				Delete(gomock.Any(), idMatcher).
				Return(c.returnedErr).
				Times(c.expectTime)

			grpc := New(m)

			res, err := grpc.DeleteCompany(context.Background(), c.request)

			if c.wantErr && err != nil {
				if c.returnedErr != nil && err == nil {
					t.Fatalf("expect err %s", err.Error())
				}
				return
			}

			if v, e := res.Id, c.request.Id; v != e {
				t.Fatalf("ecpect id: %d but %d", v, e)
			}
		})
	}
}

var testUpdateCompanyCase = []struct {
	name            string
	request         *pb.CompanyRequest
	expectTime      int
	returnedCompany *model.Company
	returnedErr     error
	wantErr         bool
}{
	{
		name: "happy",
		request: &pb.CompanyRequest{
			Id: api_model.Int64ToPointer(11),
			Company: &pb.Company{
				Name: "any",
			},
		},
		expectTime: 1,
		returnedCompany: &model.Company{
			ID:   1,
			Name: "any",
		},
		returnedErr: nil,
		wantErr:     false,
	},
	{
		name: "error - not found",
		request: &pb.CompanyRequest{
			Id:      api_model.Int64ToPointer(2211),
			Company: &pb.Company{},
		},
		expectTime:      1,
		returnedCompany: nil,
		returnedErr:     nil,
		wantErr:         true,
	},
	{
		name:            "error - no id",
		request:         &pb.CompanyRequest{},
		expectTime:      0,
		returnedCompany: nil,
		returnedErr:     nil,
		wantErr:         true,
	},
	{
		name: "error - servce return err",
		request: &pb.CompanyRequest{
			Id:      api_model.Int64ToPointer(11),
			Company: &pb.Company{},
		},
		expectTime:      1,
		returnedCompany: nil,
		returnedErr:     errors.New("any"),
		wantErr:         true,
	},
}

func TestGrpcService_UpdateCompany(t *testing.T) {
	for _, c := range testUpdateCompanyCase {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			idMatcher := gomock.Any()
			if c.request != nil && c.request.Id != nil {
				idMatcher = gomock.Eq(*c.request.Id)
			}

			m := mock.NewMockservice(ctrl)
			m.EXPECT().
				Update(gomock.Any(), idMatcher, gomock.Any()).
				Return(c.returnedCompany, c.returnedErr).
				Times(c.expectTime)

			grpc := New(m)

			res, err := grpc.UpdateCompany(context.Background(), c.request)

			if c.wantErr && err != nil {
				if c.returnedErr != nil && err == nil {
					t.Fatalf("expect err %s", err.Error())
				}
				return
			}

			if res == nil {
				t.Fatal("expect not nil")
			}

			if v, e := res.Company.Id, c.returnedCompany.ID; v != e {
				t.Fatalf("ecpect id: %d but %d", v, e)
			}
			if v, e := res.Company.Name, c.returnedCompany.Name; v != e {
				t.Fatalf("ecpect name: %s but %s", v, e)
			}
		})
	}
}

var testCreateCompanyCase = []struct {
	name        string
	request     *pb.CompanyRequest
	expectTime  int
	returnedID  int64
	returnedErr error
	wantErr     bool
}{
	{
		name: "happy",
		request: &pb.CompanyRequest{
			Company: &pb.Company{
				Name: "any",
			},
		},
		expectTime:  1,
		returnedID:  1,
		returnedErr: nil,
		wantErr:     false,
	},
	{
		name:        "error - without name",
		request:     &pb.CompanyRequest{},
		expectTime:  0,
		returnedErr: nil,
		wantErr:     true,
	},
	{
		name: "error - with id",
		request: &pb.CompanyRequest{
			Id: api_model.Int64ToPointer(11),
		},
		expectTime:  0,
		returnedErr: nil,
		wantErr:     true,
	},
	{
		name: "error - create failed",
		request: &pb.CompanyRequest{
			Company: &pb.Company{
				Name: "any",
			},
		},
		expectTime:  1,
		returnedErr: errors.New("any"),
		wantErr:     true,
	},
}

func TestGrpcService_CreateCompany(t *testing.T) {
	for _, c := range testCreateCompanyCase {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			companyMatch := gomock.Any()
			if c.request != nil && c.request.Company != nil {
				companyMatch = gomock.Eq(&api_model.UpdateCompanyParams{Name: c.request.Company.Name})
			}

			m := mock.NewMockservice(ctrl)
			m.EXPECT().
				Create(gomock.Any(), companyMatch).
				Return(c.returnedID, c.returnedErr).
				Times(c.expectTime)

			grpc := New(m)

			res, err := grpc.CreateCompany(context.Background(), c.request)

			if c.wantErr && err != nil {
				if c.returnedErr != nil && err == nil {
					t.Fatalf("expect err %s", err.Error())
				}
				return
			}

			if v, e := *res.Id, c.returnedID; v != e {
				t.Fatalf("ecpect %d but %d", v, e)
			}
		})
	}
}
