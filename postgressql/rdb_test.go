package postgressql

// Actually functional testing, using in memory database to testing.

import (
	"context"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/sean0427/company-domain-distributed-system-p/api_model"
	"github.com/sean0427/company-domain-distributed-system-p/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	testingDB *gorm.DB
)

func TestMain(m *testing.M) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	testingDB = db
	testingDB.AutoMigrate(&model.Company{})

	m.Run()
	os.Exit(0)
}

var testGetCompaniesCases = []struct {
	name       string
	testCount  int
	testParams api_model.GetCompaniesParams
	wantCount  int
	wantError  bool
}{
	{
		name:      "zero path - GetCompanies",
		testCount: 0,
		wantCount: 0,
		wantError: false,
	},
	{
		name:      "happy path - GetCompanies",
		testCount: 10,
		wantCount: 10,
		wantError: false,
	},
	{
		name:      "error path",
		wantCount: 10,
		wantError: true,
	},
	{
		name: "filter path - fullname contains 1",
		testParams: api_model.GetCompaniesParams{
			Name: api_model.StringToPointer("1"),
		},
		testCount: 20,
		wantCount: 3, // 1, 10, 11
		wantError: false,
	},
	{
		name: "filter path - fullname contains test",
		testParams: api_model.GetCompaniesParams{
			Name: api_model.StringToPointer("test"),
		},
		testCount: 20,
		wantCount: 20,
		wantError: false,
	},
}

func Test_reposity_Get(t *testing.T) {
	for _, c := range testGetCompaniesCases {

		t.Run(c.name, func(t *testing.T) {
			createRandomCompanyToDB(c.testCount)
			testParams := c.testParams
			repo := repository{db: testingDB}

			prodct, err := repo.Get(context.Background(), &testParams)

			if err != nil && !c.wantError {
				t.Errorf("got error %v", err)
				return
			}
			if len(prodct) != c.wantCount {
				t.Errorf("Expected %d companies, got %d", c.wantCount, len(prodct))
			}
		})
	}
}

var testGetCompanyIdCases = []struct {
	name      string
	id        int64
	testCount int
	want      string
	wantError bool
}{
	{
		name:      "happy - get company id",
		id:        0,
		testCount: 1,
		wantError: false,
	},
	{
		name:      "happy - get company id 2",
		id:        1,
		testCount: 2,
		wantError: false,
	},
	{
		name:      "happy - get company id 100",
		id:        99,
		testCount: 100,
		wantError: false,
	},
	{
		name:      "error - not create",
		testCount: 0,
		id:        1,
		wantError: true,
	},
}

func Test_repository_GetByID(t *testing.T) {
	for _, c := range testGetCompanyIdCases {
		t.Run(c.name, func(t *testing.T) {
			createRandomCompanyToDB(c.testCount)
			repo := repository{db: testingDB}
			ctx := context.WithValue(context.Background(), api_model.CtxKey("user"), "aa")
			company, err := repo.GetByID(ctx, c.id)

			if err != nil && !c.wantError {
				t.Errorf("got error %v", err)
				return
			}

			if company.ID == c.id {
				t.Errorf("Expected %d, got %d", c.id, company.ID)
			}
		})
	}
}

func CtxKey(s string) {
	panic("unimplemented")
}

func Test_repository_Create(t *testing.T) {
	tests := []struct {
		name    string
		params  *api_model.CreateCompanyParams
		want    int64
		wantErr bool
	}{
		{
			name: "happy",
			params: &api_model.CreateCompanyParams{
				Name:    "test",
				Email:   "test@test.com",
				Contact: "featea",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "error",
			params: &api_model.CreateCompanyParams{
				Name: "",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := repository{db: testingDB}
			ctx := context.WithValue(context.Background(), api_model.CtxKey("user"), "aa")

			got, err := r.Create(ctx, tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("repository.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repository_Delete(t *testing.T) {
	createRandomCompanyToDB(10)

	tests := []struct {
		name    string
		id      int64
		wantErr bool
	}{
		{
			name:    "happy",
			id:      1,
			wantErr: false,
		},
		{
			name:    "error",
			id:      11,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{db: testingDB}

			ctx := context.WithValue(context.Background(), api_model.CtxKey("user"), "aa")
			if err := r.Delete(ctx, tt.id); (err != nil) != tt.wantErr {
				t.Errorf("repository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_repository_Update(t *testing.T) {
	createRandomCompanyToDB(1)
	type args struct {
		id     int64
		params *api_model.UpdateCompanyParams
	}
	tests := []struct {
		name    string
		args    args
		want    *model.Company
		wantErr bool
	}{
		{
			name: "happy",
			args: args{
				id: 1,
				params: &api_model.UpdateCompanyParams{
					ID:   1,
					Name: "test"},
			},
			want: &model.Company{
				ID:   1,
				Name: "test",
			},
			wantErr: false,
		}, {
			name: "error",

			args: args{
				id: 100,
				params: &api_model.UpdateCompanyParams{
					ID:   100,
					Name: "test"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &repository{
				db: testingDB,
			}
			got, err := r.Update(context.Background(), tt.args.id, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("repository.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("repository.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func createRandomCompanyToDB(numbers int) {
	for i := 0; i < numbers; i++ {
		company := &model.Company{
			ID:        int64(i),
			Name:      "test" + strconv.Itoa(i),
			CreatedBy: strconv.Itoa(i + 1),
		}
		testingDB.Create(company)
	}
}
