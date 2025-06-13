package service

import (
	"github.com/ronannnn/infra/model"
	"gorm.io/gorm"
)

type ApiRepo interface {
	CrudRepo[model.Api]
}

func NewApiService(
	repo ApiRepo,
	db *gorm.DB,
) *ApiService {
	return &ApiService{
		DefaultCrudSrv[model.Api]{
			Repo: repo,
			Db:   db,
		},
	}
}

type ApiService struct {
	DefaultCrudSrv[model.Api]
}
