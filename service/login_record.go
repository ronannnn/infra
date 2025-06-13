package service

import (
	"github.com/ronannnn/infra/model"
	"gorm.io/gorm"
)

type LoginRecordRepo interface {
	CrudRepo[model.LoginRecord]
}

func NewLoginRecordService(
	repo LoginRecordRepo,
	db *gorm.DB,
) *LoginRecordService {
	return &LoginRecordService{
		DefaultCrudSrv[model.LoginRecord]{
			Repo: repo,
			Db:   db,
		},
	}
}

type LoginRecordService struct {
	DefaultCrudSrv[model.LoginRecord]
}
