package service

import (
	"github.com/ronannnn/infra/model"
	"gorm.io/gorm"
)

type DepartmentRepo interface {
	CrudRepo[model.Department]
}

func NewDepartmentService(
	repo DepartmentRepo,
	db *gorm.DB,
) *DepartmentService {
	return &DepartmentService{
		DefaultCrudSrv[model.Department]{
			Repo: repo,
			Db:   db,
		},
	}
}

type DepartmentService struct {
	DefaultCrudSrv[model.Department]
}
