package api

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Store interface {
	create(model *models.Api) error
	update(partialUpdatedModel *models.Api) (models.Api, error)
	deleteById(id uint) error
	deleteByIds(ids []uint) error
	list(query query.ApiQuery) (response.PageResult, error)
	getById(id uint) (models.Api, error)
}

func ProvideStore(db *gorm.DB) Store {
	return StoreImpl{db: db}
}

type StoreImpl struct {
	db *gorm.DB
}

func (s StoreImpl) create(model *models.Api) error {
	return s.db.Create(model).Error
}

func (s StoreImpl) update(partialUpdatedModel *models.Api) (updatedModel models.Api, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	result := s.db.Updates(partialUpdatedModel)
	if result.Error != nil {
		return updatedModel, result.Error
	}
	if result.RowsAffected == 0 {
		return updatedModel, models.ErrModified("Api")
	}
	return s.getById(partialUpdatedModel.Id)
}

func (s StoreImpl) deleteById(id uint) error {
	return s.db.Delete(&models.Api{}, "id = ?", id).Error
}

func (s StoreImpl) deleteByIds(ids []uint) error {
	return s.db.Delete(&models.Api{}, "id IN ?", ids).Error
}

func (s StoreImpl) list(apiQuery query.ApiQuery) (result response.PageResult, err error) {
	var total int64
	var list []models.Api
	if err = s.db.Model(&models.Api{}).Count(&total).Error; err != nil {
		return
	}
	if err = s.db.
		Scopes(query.MakeConditionFromQuery(apiQuery)).
		Scopes(query.Paginate(apiQuery.Pagination.PageNum, apiQuery.Pagination.PageSize)).
		Find(&list).Error; err != nil {
		return
	}
	result = response.PageResult{
		List:     list,
		Total:    total,
		PageNum:  apiQuery.Pagination.PageNum,
		PageSize: apiQuery.Pagination.PageSize,
	}
	return
}

func (s StoreImpl) getById(id uint) (model models.Api, err error) {
	err = s.db.First(&model, "id = ?", id).Error
	return
}
