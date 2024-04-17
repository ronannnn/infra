package menu

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Store interface {
	Create(tx *gorm.DB, model *Menu) error
	Update(tx *gorm.DB, partialUpdatedModel *Menu) (Menu, error)
	DeleteById(tx *gorm.DB, id uint) error
	DeleteByIds(tx *gorm.DB, ids []uint) error
	List(tx *gorm.DB, query MenuQuery) (response.PageResult, error)
	GetById(tx *gorm.DB, id uint) (Menu, error)
}

func ProvideStore() Store {
	return StoreImpl{}
}

type StoreImpl struct {
}

func (s StoreImpl) Create(tx *gorm.DB, model *Menu) error {
	return tx.Create(model).Error
}

func (s StoreImpl) Update(tx *gorm.DB, partialUpdatedModel *Menu) (updatedModel Menu, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	if err = tx.Transaction(func(tx *gorm.DB) (err error) {
		// update associations with Associations()
		if partialUpdatedModel.Apis != nil {
			if err = tx.Model(partialUpdatedModel).Association("Apis").Replace(partialUpdatedModel.Apis); err != nil {
				return err
			}
			// set associations to nil to avoid Updates() below,
			partialUpdatedModel.Apis = nil
		}
		// update all other non-associations
		// if no other fields are updated, it still update the version so no error will occur
		result := tx.Updates(partialUpdatedModel)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return models.ErrModified("Menu")
		}
		return
	}); err != nil {
		return updatedModel, err
	}
	return s.GetById(tx, partialUpdatedModel.Id)
}

func (s StoreImpl) DeleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&Menu{}, "id = ?", id).Error
}

func (s StoreImpl) DeleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&Menu{}, "id IN ?", ids).Error
}

func (s StoreImpl) List(tx *gorm.DB, menuQuery MenuQuery) (result response.PageResult, err error) {
	var total int64
	var list []Menu
	if err = tx.Model(&Menu{}).Count(&total).Error; err != nil {
		return
	}
	if err = tx.
		Scopes(query.MakeConditionFromQuery(menuQuery)).
		Preload("Apis").
		Find(&list).Error; err != nil {
		return
	}
	result = response.PageResult{
		List:     list,
		Total:    total,
		PageNum:  1,
		PageSize: int(total),
	}
	return
}

func (s StoreImpl) GetById(tx *gorm.DB, id uint) (model Menu, err error) {
	err = tx.Preload("Apis").First(&model, "id = ?", id).Error
	return
}
