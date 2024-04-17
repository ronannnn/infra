package role

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Store interface {
	Create(tx *gorm.DB, model *Role) error
	Update(tx *gorm.DB, partialUpdatedModel *Role) (Role, error)
	DeleteById(tx *gorm.DB, id uint) error
	DeleteByIds(tx *gorm.DB, ids []uint) error
	List(tx *gorm.DB, query RoleQuery) (response.PageResult, error)
	GetById(tx *gorm.DB, id uint) (Role, error)
}

func ProvideStore() Store {
	return StoreImpl{}
}

type StoreImpl struct {
}

func (s StoreImpl) Create(tx *gorm.DB, model *Role) error {
	return tx.Create(model).Error
}

func (s StoreImpl) Update(tx *gorm.DB, partialUpdatedModel *Role) (updatedModel Role, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	if err = tx.Transaction(func(tx *gorm.DB) (err error) {
		// update associations with Associations()
		if partialUpdatedModel.Menus != nil {
			if err = tx.Model(partialUpdatedModel).Association("Menus").Replace(partialUpdatedModel.Menus); err != nil {
				return err
			}
			// set associations to nil to avoid Updates() below,
			partialUpdatedModel.Menus = nil
		}
		// update all other non-associations
		// if no other fields are updated, it still update the version so no error will occur
		result := tx.Updates(partialUpdatedModel)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return models.ErrModified("Role")
		}
		return
	}); err != nil {
		return updatedModel, err
	}
	return s.GetById(tx, partialUpdatedModel.Id)
}

func (s StoreImpl) DeleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&Role{}, "id = ?", id).Error
}

func (s StoreImpl) DeleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&Role{}, "id IN ?", ids).Error
}

func (s StoreImpl) List(tx *gorm.DB, roleQuery RoleQuery) (result response.PageResult, err error) {
	var total int64
	var list []Role
	if err = tx.Model(&Role{}).Count(&total).Error; err != nil {
		return
	}
	if err = tx.
		Scopes(query.MakeConditionFromQuery(roleQuery)).
		Scopes(query.Paginate(roleQuery.Pagination.PageNum, roleQuery.Pagination.PageSize)).
		Preload("Menus").
		Find(&list).Error; err != nil {
		return
	}
	result = response.PageResult{
		List:     list,
		Total:    total,
		PageNum:  roleQuery.Pagination.PageNum,
		PageSize: roleQuery.Pagination.PageSize,
	}
	return
}

func (s StoreImpl) GetById(tx *gorm.DB, id uint) (model Role, err error) {
	err = tx.Preload("Menus").First(&model, "id = ?", id).Error
	return
}
