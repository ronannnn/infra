package repo

import (
	"context"

	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/model/request/query"
	"github.com/ronannnn/infra/model/response"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"gorm.io/gorm"
)

// DefaultCrudRepo is a default implementation of the Services.Crud interface
type DefaultCrudRepo[T model.Crudable] struct {
	preloads []string // 用于gorm的Preload
}

func NewDefaultCrudRepo[T model.Crudable](
	preloads ...string,
) DefaultCrudRepo[T] {
	return DefaultCrudRepo[T]{
		preloads: preloads,
	}
}

func (crud DefaultCrudRepo[T]) Create(ctx context.Context, tx *gorm.DB, model *T) error {
	if model == nil {
		return msg.NewError(reason.DbModelCreatedError).WithStack()
	}
	return tx.WithContext(ctx).Create(model).Error
}

func (crud DefaultCrudRepo[T]) CreateWithScopes(ctx context.Context, tx *gorm.DB, model *T) (updatedModel *T, err error) {
	if model == nil {
		err = msg.NewError(reason.DbModelCreatedError).WithStack()
		return
	}
	if err = tx.WithContext(ctx).Create(&model).Error; err != nil {
		return
	}
	return crud.GetById(ctx, tx, (*model).GetId())
}

func (crud DefaultCrudRepo[T]) Update(ctx context.Context, tx *gorm.DB, partialUpdatedModel *T) (updatedModel *T, err error) {
	if partialUpdatedModel == nil {
		return nil, msg.NewError(reason.DbModelUpdatedError).WithStack()
	}
	if (*partialUpdatedModel).GetId() == 0 {
		return updatedModel, msg.NewError(reason.DbModelUpdatedIdCannotBeZero).WithStack()
	}
	err = tx.WithContext(ctx).Transaction(func(tx2 *gorm.DB) (err error) {
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
		updatedModel, err = crud.GetById(ctx, tx2, (*partialUpdatedModel).GetId())
		return
	})
	return
}

func (crud DefaultCrudRepo[T]) DeleteById(ctx context.Context, tx *gorm.DB, id uint) error {
	return tx.WithContext(ctx).Transaction(func(tx2 *gorm.DB) (err error) {
		var t T
		if err = tx2.Delete(&t, "id = ?", id).Error; err != nil {
			return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
		}
		return
	})
}

func (crud DefaultCrudRepo[T]) DeleteByIds(ctx context.Context, tx *gorm.DB, ids []uint) (err error) {
	var t T
	if err = tx.WithContext(ctx).WithContext(ctx).Delete(&t, "id IN ?", ids).Error; err != nil {
		return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (crud DefaultCrudRepo[T]) List(ctx context.Context, tx *gorm.DB, apiQuery query.Query) (result *response.PageResult[T], err error) {
	txWithCtx := tx.WithContext(ctx)
	var t T
	var total int64
	var list []*T
	if err = txWithCtx.Model(&t).Count(&total).Error; err != nil {
		return
	}
	queryScope, err := query.MakeConditionFromQuery(apiQuery, t)
	if err != nil {
		return nil, msg.NewError(reason.DatabaseError).WithError(err).WithStack()
	}
	if err = txWithCtx.
		Scopes(crud.Preload()).
		Scopes(queryScope).
		Scopes(query.Paginate(apiQuery.Pagination)).
		Find(&list).Error; err != nil {
		return nil, msg.NewError(reason.DatabaseError).WithError(err).WithStack()
	}
	result = &response.PageResult[T]{
		List:     list,
		Total:    total,
		PageNum:  apiQuery.Pagination.PageNum,
		PageSize: apiQuery.Pagination.PageSize,
	}
	return
}

func (crud DefaultCrudRepo[T]) GetById(ctx context.Context, tx *gorm.DB, id uint) (model *T, err error) {
	if err = tx.
		WithContext(ctx).
		Scopes(crud.Preload()).
		First(&model, "id = ?", id).
		Error; err == gorm.ErrRecordNotFound {
		err = msg.NewError(reason.DbModelReadIdNotExists).WithReasonTemplateData(reason.IdTd{Id: id}).WithStack()
		return
	} else if err != nil {
		err = msg.NewError(reason.DatabaseError).WithError(err).WithStack()
		return
	}
	return
}

func (crud DefaultCrudRepo[T]) Preload() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, preload := range crud.preloads {
			db = db.Preload(preload)
		}
		return db
	}
}
