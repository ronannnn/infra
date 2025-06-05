package service

import (
	"github.com/ronannnn/infra/model"
	"gorm.io/gorm"
)

type RowRecordRepo interface {
	CrudRepo[*model.RowRecord]
	GetByTableNameAndRowId(tx *gorm.DB, tableName string, rowId uint) (list []model.RowRecord, err error)
}

func NewRowRecordService(
	repo RowRecordRepo,
	db *gorm.DB,
) *RowRecordService {
	return &RowRecordService{
		DefaultCrudSrv[*model.RowRecord]{
			Repo: repo,
			Db:   db,
		},
	}
}

type RowRecordService struct {
	DefaultCrudSrv[*model.RowRecord]
}
