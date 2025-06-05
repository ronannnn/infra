package service

import (
	"github.com/ronannnn/infra/model"
	"gorm.io/gorm"
)

type ApiRecordRepo interface {
	CrudRepo[*model.ApiRecord]
}

func NewApiRecordService(
	repo ApiRecordRepo,
	db *gorm.DB,
) *ApiRecordService {
	return &ApiRecordService{
		DefaultCrudSrv[*model.ApiRecord]{
			Repo: repo,
			Db:   db,
		},
	}
}

type ApiRecordService struct {
	DefaultCrudSrv[*model.ApiRecord]
}
