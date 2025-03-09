package repos

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"gorm.io/gorm"
)

// DefaultCrudRepo is a default implementation of the Services.Crud interface
type DefaultCrudRepo[T models.Crudable] struct{}

func NewDefaultCrudRepo[T models.Crudable]() DefaultCrudRepo[T] {
	return DefaultCrudRepo[T]{}
}

func (crud DefaultCrudRepo[T]) Create(tx *gorm.DB, model *T) error {
	return tx.Create(model).Error
}

func (crud DefaultCrudRepo[T]) CreateWithScopes(tx *gorm.DB, model *T) (updatedModel *T, err error) {
	if err = tx.Create(model).Error; err != nil {
		return
	}
	return crud.GetById(tx, (*model).GetId())
}

func (crud DefaultCrudRepo[T]) Update(tx *gorm.DB, partialUpdatedModel *T) (updatedModel *T, err error) {
	if partialUpdatedModel == nil || (*partialUpdatedModel).GetId() == 0 {
		return updatedModel, msg.NewError(reason.DbModelUpdatedIdCannotBeZero).WithStack()
	}
	err = tx.Transaction(func(tx2 *gorm.DB) (err error) {
		// update associations with Associations()
		// ...
		// update all other non-associations
		// if no other fields are updated, it still update the version so no error will occur
		result := tx2.Updates(partialUpdatedModel)
		if result.Error != nil {
			return msg.NewError(reason.DatabaseError).WithError(result.Error).WithStack()
		}
		if result.RowsAffected == 0 {
			return msg.NewError(reason.DbModelAlreadyUpdatedByOthers).WithStack()
		}
		updatedModel, err = crud.GetById(tx2, (*partialUpdatedModel).GetId())
		return
	})
	return
}

func (crud DefaultCrudRepo[T]) DeleteById(tx *gorm.DB, id uint) error {
	var t T
	return tx.Delete(&t, "id = ?", id).Error
}

func (crud DefaultCrudRepo[T]) DeleteByIds(tx *gorm.DB, ids []uint) error {
	var t T
	return tx.Delete(&t, "id IN ?", ids).Error
}

func (crud DefaultCrudRepo[T]) List(tx *gorm.DB, apiQuery query.Query) (result response.PageResult, err error) {
	var t T
	var total int64
	var list []T
	if err = tx.Model(&t).Count(&total).Error; err != nil {
		return
	}
	queryScope, err := query.MakeConditionFromQuery(apiQuery, t)
	if err != nil {
		return
	}
	if err = tx.
		Scopes(crud.Preload()).
		Scopes(queryScope).
		Scopes(query.Paginate(apiQuery.Pagination)).
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

func (crud DefaultCrudRepo[T]) GetById(tx *gorm.DB, id uint) (model *T, err error) {
	err = tx.Scopes(crud.Preload()).First(model, "id = ?", id).Error
	return
}

func (crud DefaultCrudRepo[T]) Preload() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}
