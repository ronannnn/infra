package service

import (
	"github.com/ronannnn/infra/model"
	"gorm.io/gorm"
)

type RoleRepo interface {
	CrudRepo[model.Role]
}

func NewRoleService(
	repo RoleRepo,
	db *gorm.DB,
) *RoleService {
	return &RoleService{
		DefaultCrudSrv[model.Role]{
			Repo: repo,
			Db:   db,
		},
	}
}

type RoleService struct {
	DefaultCrudSrv[model.Role]
}
