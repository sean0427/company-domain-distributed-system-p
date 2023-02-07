package api_model

type GetCompaniesParams struct {
	Name *string
}

type CreateCompanyParams struct {
	Name     string
	Email    string
	Password string
}

type UpdateCompanyParams struct {
	ID       int64
	Name     string
	Email    string
	Password string
}

func StringToPointer(s string) *string {
	return &s
}

func Int64ToPointer(i int64) *int64 {
	return &i
}