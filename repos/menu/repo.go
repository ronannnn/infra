package menu

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	srv "github.com/ronannnn/infra/services/menu"
	"gorm.io/gorm"
)

func New() srv.Repo {
	return &repo{}
}

type repo struct {
}

func (s repo) Create(tx *gorm.DB, model *models.Menu) error {
	return tx.Create(model).Error
}

func (s repo) Update(tx *gorm.DB, partialUpdatedModel *models.Menu) (updatedModel models.Menu, err error) {
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
			return models.ErrModified("models.Menu")
		}
		return
	}); err != nil {
		return updatedModel, err
	}
	return s.GetById(tx, partialUpdatedModel.Id)
}

func (s repo) DeleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.Menu{}, "id = ?", id).Error
}

func (s repo) DeleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&models.Menu{}, "id IN ?", ids).Error
}

func (s repo) List(tx *gorm.DB, menuQuery query.Query) (result response.PageResult, err error) {
	var total int64
	var list []models.Menu
	if err = tx.Model(&models.Menu{}).Count(&total).Error; err != nil {
		return
	}
	queryScope, err := query.MakeConditionFromQuery(menuQuery, models.Menu{})
	if err != nil {
		return
	}
	if err = tx.
		Scopes(queryScope).
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

func (s repo) GetById(tx *gorm.DB, id uint) (model models.Menu, err error) {
	err = tx.Preload("Apis").First(&model, "id = ?", id).Error
	return
}
