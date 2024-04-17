package department

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Store interface {
	Create(tx *gorm.DB, model *Department) error
	Update(tx *gorm.DB, partialUpdatedModel *Department) (Department, error)
	DeleteById(tx *gorm.DB, id uint) error
	DeleteByIds(tx *gorm.DB, ids []uint) error
	List(tx *gorm.DB, query DepartmentQuery) (response.PageResult, error)
	GetById(tx *gorm.DB, id uint) (Department, error)
}

func ProvideStore() Store {
	return StoreImpl{}
}

type StoreImpl struct {
}

func (s StoreImpl) Create(tx *gorm.DB, model *Department) error {
	return tx.Create(model).Error
}

func (s StoreImpl) Update(tx *gorm.DB, partialUpdatedModel *Department) (updatedModel Department, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	if err = tx.Updates(partialUpdatedModel).Error; err != nil {
		return
	}
	return s.GetById(tx, partialUpdatedModel.Id)
}

func (s StoreImpl) DeleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&Department{}, "id = ?", id).Error
}

func (s StoreImpl) DeleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&Department{}, "id IN ?", ids).Error
}

func (s StoreImpl) List(tx *gorm.DB, menuQuery DepartmentQuery) (result response.PageResult, err error) {
	var total int64
	var list []Department
	if err = tx.Model(&Department{}).Count(&total).Error; err != nil {
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

func (s StoreImpl) GetById(tx *gorm.DB, id uint) (model Department, err error) {
	err = tx.Preload("Apis").First(&model, "id = ?", id).Error
	return
}
