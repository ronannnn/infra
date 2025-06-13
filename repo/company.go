package repo

import (
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewCompanyRepo() service.CompanyRepo {
	return &companyRepo{}
}

type companyRepo struct {
	DefaultCrudRepo[model.Company]
}
