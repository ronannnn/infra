package company

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	srv "github.com/ronannnn/infra/services/company"
	"gorm.io/gorm"
)

func New() srv.Repo {
	return &repo{}
}

type repo struct {
}

func (r repo) Create(tx *gorm.DB, model *models.Company) error {
	return tx.Create(model).Error
}

func (r repo) Update(tx *gorm.DB, partialUpdatedModel *models.Company) (updatedModel *models.Company, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	if err = tx.Updates(partialUpdatedModel).Error; err != nil {
		return
	}
	return r.GetById(tx, partialUpdatedModel.Id)
}

func (r repo) DeleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.Company{}, "id = ?", id).Error
}

func (r repo) DeleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&models.Company{}, "id IN ?", ids).Error
}

func (r repo) List(tx *gorm.DB, menuQuery query.Query) (result response.PageResult, err error) {
	var total int64
	var list []models.Company
	if err = tx.Model(&models.Company{}).Count(&total).Error; err != nil {
		return
	}
	queryScope, err := query.MakeConditionFromQuery(menuQuery, models.Company{})
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

func (r repo) GetById(tx *gorm.DB, id uint) (model *models.Company, err error) {
	err = tx.Preload("Apis").First(&model, "id = ?", id).Error
	return
}
