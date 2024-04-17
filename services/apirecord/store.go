package apirecord

import (
	"github.com/ronannnn/infra/models"
	"gorm.io/gorm"
)

type Store interface {
	create(*gorm.DB, *ApiRecord) error
	update(*gorm.DB, *ApiRecord) (ApiRecord, error)
	deleteById(*gorm.DB, uint) error
	deleteByIds(*gorm.DB, []uint) error
	getById(*gorm.DB, uint) (ApiRecord, error)
}

type StoreImpl struct {
}

func ProvideStore(db *gorm.DB) Store {
	return &StoreImpl{}
}

func (s StoreImpl) create(tx *gorm.DB, model *ApiRecord) error {
	return tx.Create(model).Error
}

func (s StoreImpl) update(tx *gorm.DB, partialUpdatedModel *ApiRecord) (updatedModel ApiRecord, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	result := tx.Updates(partialUpdatedModel)
	if result.Error != nil {
		return updatedModel, result.Error
	}
	if result.RowsAffected == 0 {
		return updatedModel, models.ErrModified("ApiRecord")
	}
	return s.getById(tx, partialUpdatedModel.Id)
}

func (s StoreImpl) deleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&ApiRecord{}, "id = ?", id).Error
}

func (s StoreImpl) deleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&ApiRecord{}, "id IN ?", ids).Error
}

func (s StoreImpl) getById(tx *gorm.DB, id uint) (model ApiRecord, err error) {
	err = tx.First(&model, "id = ?", id).Error
	return
}
