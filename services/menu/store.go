package menu

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Store interface {
	create(model *models.Menu) error
	update(partialUpdatedModel *models.Menu) (models.Menu, error)
	deleteById(id uint) error
	deleteByIds(ids []uint) error
	list(query query.MenuQuery) (response.PageResult, error)
	getById(id uint) (models.Menu, error)
}

func ProvideStore(db *gorm.DB) Store {
	return StoreImpl{db: db}
}

type StoreImpl struct {
	db *gorm.DB
}

func (s StoreImpl) create(model *models.Menu) error {
	return s.db.Create(model).Error
}

func (s StoreImpl) update(partialUpdatedModel *models.Menu) (updatedModel models.Menu, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	if err = s.db.Transaction(func(tx *gorm.DB) (err error) {
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
	return s.getById(partialUpdatedModel.Id)
}

func (s StoreImpl) deleteById(id uint) error {
	return s.db.Delete(&models.Menu{}, "id = ?", id).Error
}

func (s StoreImpl) deleteByIds(ids []uint) error {
	return s.db.Delete(&models.Menu{}, "id IN ?", ids).Error
}

func (s StoreImpl) list(menuQuery query.MenuQuery) (result response.PageResult, err error) {
	var total int64
	var list []models.Menu
	if err = s.db.Model(&models.Menu{}).Count(&total).Error; err != nil {
		return
	}
	if err = s.db.
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

func (s StoreImpl) getById(id uint) (model models.Menu, err error) {
	err = s.db.Preload("Apis").First(&model, "id = ?", id).Error
	return
}
