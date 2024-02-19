package apirecord

import (
	"github.com/ronannnn/infra/models"
	"gorm.io/gorm"
)

type Store interface {
	create(model *models.ApiRecord) error
	update(partialUpdatedModel *models.ApiRecord) (models.ApiRecord, error)
	deleteById(id uint) error
	deleteByIds(ids []uint) error
	getById(id uint) (models.ApiRecord, error)
}

type StoreImpl struct {
	db *gorm.DB
}

func ProvideStore(db *gorm.DB) Store {
	return &StoreImpl{db: db}
}

func (s StoreImpl) create(model *models.ApiRecord) error {
	return s.db.Create(model).Error
}

func (s StoreImpl) update(partialUpdatedModel *models.ApiRecord) (updatedModel models.ApiRecord, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	result := s.db.Updates(partialUpdatedModel)
	if result.Error != nil {
		return updatedModel, result.Error
	}
	if result.RowsAffected == 0 {
		return updatedModel, models.ErrModified("ApiRecord")
	}
	return s.getById(partialUpdatedModel.Id)
}

func (s StoreImpl) deleteById(id uint) error {
	return s.db.Delete(&models.ApiRecord{}, "id = ?", id).Error
}

func (s StoreImpl) deleteByIds(ids []uint) error {
	return s.db.Delete(&models.ApiRecord{}, "id IN ?", ids).Error
}

func (s StoreImpl) getById(id uint) (model models.ApiRecord, err error) {
	err = s.db.First(&model, "id = ?", id).Error
	return
}
