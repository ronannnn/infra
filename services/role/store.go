package role

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Store interface {
	create(model *models.Role) error
	update(partialUpdatedModel *models.Role) (models.Role, error)
	deleteById(id uint) error
	deleteByIds(ids []uint) error
	list(query query.RoleQuery) (response.PageResult, error)
	getById(id uint) (models.Role, error)
}

func ProvideStore(db *gorm.DB) Store {
	return StoreImpl{db: db}
}

type StoreImpl struct {
	db *gorm.DB
}

func (s StoreImpl) create(model *models.Role) error {
	return s.db.Create(model).Error
}

func (s StoreImpl) update(partialUpdatedModel *models.Role) (updatedModel models.Role, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	if err = s.db.Transaction(func(tx *gorm.DB) (err error) {
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
	return s.getById(partialUpdatedModel.Id)
}

func (s StoreImpl) deleteById(id uint) error {
	return s.db.Delete(&models.Role{}, "id = ?", id).Error
}

func (s StoreImpl) deleteByIds(ids []uint) error {
	return s.db.Delete(&models.Role{}, "id IN ?", ids).Error
}

func (s StoreImpl) list(roleQuery query.RoleQuery) (result response.PageResult, err error) {
	var total int64
	var list []models.Role
	if err = s.db.Model(&models.Role{}).Count(&total).Error; err != nil {
		return
	}
	if err = s.db.
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

func (s StoreImpl) getById(id uint) (model models.Role, err error) {
	err = s.db.Preload("Menus").First(&model, "id = ?", id).Error
	return
}
