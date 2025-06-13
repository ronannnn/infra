package service

import (
	"github.com/ronannnn/infra/model"
	"gorm.io/gorm"
)

type MenuRepo interface {
	CrudRepo[model.Menu]
}

func NewMenuService(
	repo MenuRepo,
	db *gorm.DB,
) *MenuService {
	return &MenuService{
		DefaultCrudSrv[model.Menu]{
			Repo: repo,
			Db:   db,
		},
	}
}

type MenuService struct {
	DefaultCrudSrv[model.Menu]
}
