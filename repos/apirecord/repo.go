package apirecord

import (
	"github.com/ronannnn/infra/models"
	srv "github.com/ronannnn/infra/services/apirecord"
	"gorm.io/gorm"
)

type repo struct {
}

func ProvideStore(db *gorm.DB) srv.Repo {
	return &repo{}
}

func (r repo) Create(tx *gorm.DB, model *models.ApiRecord) error {
	return tx.Create(model).Error
}

func (r repo) Update(tx *gorm.DB, partialUpdatedModel *models.ApiRecord) (updatedModel *models.ApiRecord, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	result := tx.Updates(partialUpdatedModel)
	if result.Error != nil {
		return updatedModel, result.Error
	}
	if result.RowsAffected == 0 {
		return updatedModel, models.ErrModified("models.ApiRecord")
	}
	return r.GetById(tx, partialUpdatedModel.Id)
}

func (r repo) DeleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.ApiRecord{}, "id = ?", id).Error
}

func (r repo) DeleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&models.ApiRecord{}, "id IN ?", ids).Error
}

func (r repo) GetById(tx *gorm.DB, id uint) (model *models.ApiRecord, err error) {
	err = tx.First(&model, "id = ?", id).Error
	return
}
