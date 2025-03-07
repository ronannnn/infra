package role

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	srv "github.com/ronannnn/infra/services/role"
	"gorm.io/gorm"
)

func New() srv.Repo {
	return &repo{}
}

type repo struct {
}

func (s repo) Create(tx *gorm.DB, model *models.Role) error {
	return tx.Create(model).Error
}

func (s repo) Update(tx *gorm.DB, partialUpdatedModel *models.Role) (updatedModel models.Role, err error) {
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
			return models.ErrModified("models.Role")
		}
		return
	}); err != nil {
		return updatedModel, err
	}
	return s.GetById(tx, partialUpdatedModel.Id)
}

func (s repo) DeleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.Role{}, "id = ?", id).Error
}

func (s repo) DeleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&models.Role{}, "id IN ?", ids).Error
}

func (s repo) List(tx *gorm.DB, roleQuery query.Query) (result response.PageResult, err error) {
	var total int64
	var list []models.Role
	if err = tx.Model(&models.Role{}).Count(&total).Error; err != nil {
		return
	}
	queryScope, err := query.MakeConditionFromQuery(roleQuery, models.Role{})
	if err != nil {
		return
	}
	if err = tx.
		Scopes(queryScope).
		Scopes(query.Paginate(roleQuery.Pagination)).
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

func (s repo) GetById(tx *gorm.DB, id uint) (model models.Role, err error) {
	err = tx.Preload("Menus").First(&model, "id = ?", id).Error
	return
}
