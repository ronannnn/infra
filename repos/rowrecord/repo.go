package rowrecord

import (
	"github.com/ronannnn/infra/models"
	srv "github.com/ronannnn/infra/services/rowrecord"
	"gorm.io/gorm"
)

func New() srv.Repo {
	return &repo{}
}

type repo struct{}

func (s *repo) Create(tx *gorm.DB, models []*models.RowRecord) error {
	return tx.Create(&models).Error
}

func (s *repo) DeleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.RowRecord{}, "id = ?", id).Error
}

func (s *repo) DeleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&models.RowRecord{}, "id IN ?", ids).Error
}

func (s *repo) GetByTableNameAndRowId(tx *gorm.DB, tableName string, rowId uint) (list []models.RowRecord, err error) {
	err = tx.Where(&models.RowRecord{TableName: tableName, RowId: rowId}).Find(&list).Error
	return
}
