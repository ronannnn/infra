package rowrecord

import (
	"github.com/ronannnn/infra/models"
	"gorm.io/gorm"
)

type Store interface {
	Create(*gorm.DB, []models.RowRecord) error
	DeleteById(*gorm.DB, uint) error
	DeleteByIds(*gorm.DB, []uint) error
	GetByTableNameAndRowId(tx *gorm.DB, tableName string, rowId uint) ([]models.RowRecord, error)
}

func ProvideStore(db *gorm.DB) Store {
	return &StoreImpl{}
}

type StoreImpl struct{}

func (s *StoreImpl) Create(tx *gorm.DB, models []models.RowRecord) error {
	return tx.Create(&models).Error
}

func (s *StoreImpl) DeleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.RowRecord{}, "id = ?", id).Error
}

func (s *StoreImpl) DeleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&models.RowRecord{}, "id IN ?", ids).Error
}

func (s *StoreImpl) GetByTableNameAndRowId(tx *gorm.DB, tableName string, rowId uint) (list []models.RowRecord, err error) {
	err = tx.Where(&models.RowRecord{TableName: tableName, RowId: rowId}).Find(&list).Error
	return
}
