package repo

import (
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
	"gorm.io/gorm"
)

func NewRowRecordRepo() service.RowRecordRepo {
	return &rowRecordRepo{}
}

type rowRecordRepo struct {
	DefaultCrudRepo[model.RowRecord]
}

func (s *rowRecordRepo) GetByTableNameAndRowId(tx *gorm.DB, tableName string, rowId uint) (list []model.RowRecord, err error) {
	err = tx.Where(&model.RowRecord{Table: &tableName, RowId: &rowId}).Find(&list).Error
	return
}
