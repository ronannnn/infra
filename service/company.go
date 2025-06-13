package service

import (
	"github.com/ronannnn/infra/model"
	"gorm.io/gorm"
)

type CompanyRepo interface {
	CrudRepo[model.Company]
}

func NewCompanyService(
	repo CompanyRepo,
	db *gorm.DB,
) *CompanyService {
	return &CompanyService{
		DefaultCrudSrv[model.Company]{
			Repo: repo,
			Db:   db,
		},
	}
}

type CompanyService struct {
	DefaultCrudSrv[model.Company]
}
