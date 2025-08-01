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

func (crud DefaultCrudRepo[T]) CreateWithScopes(ctx context.Context, tx *gorm.DB, newModel *T) (err error) {
	if newModel == nil {
		return msg.NewError(reason.DbModelCreatedError).WithStack()
	}
	if err = tx.WithContext(ctx).Create(newModel).Error; err != nil {
		return
	}
	if err = tx.WithContext(ctx).Scopes(crud.Preload()).First(newModel, (*newModel).GetId()).Error; err != nil {
		return err
	}
	return
}

func (crud DefaultCrudRepo[T]) BatchCreate(ctx context.Context, tx *gorm.DB, models []*T) error {
	return tx.WithContext(ctx).Create(models).Error
}

func (crud DefaultCrudRepo[T]) BatchCreateWithScopes(ctx context.Context, tx *gorm.DB, models []*T) (err error) {
	if err = tx.WithContext(ctx).Create(models).Error; err != nil {
		return
	}
	ids := make([]uint, len(models))
	for i, model := range models {
		ids[i] = (*model).GetId()
	}
	if err = tx.
		Scopes(crud.Preload()).
		Find(&models, ids).Error; err != nil {
		return
	}
	return
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

func (crud DefaultCrudRepo[T]) BatchUpdate(ctx context.Context, tx *gorm.DB, partialUpdatedModels []*T) (updatedModel []*T, err error) {
	err = tx.Transaction(func(tx2 *gorm.DB) (err error) {
		updatedModel = make([]*T, len(partialUpdatedModels))
		for i := range partialUpdatedModels {
			if updatedModel[i], err = crud.Update(ctx, tx2, partialUpdatedModels[i]); err != nil {
				return
			}
		}
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

	// count
	var t T
	whereQueryScope, err := query.MakeConditionFromQuery(apiQuery, t, query.ConditionFilter{EnableWhere: true})
	if err != nil {
		return nil, msg.NewError(reason.DatabaseError).WithError(err).WithStack()
	}
	var total int64
	if err = txWithCtx.Model(&t).Scopes(whereQueryScope).Count(&total).Error; err != nil {
		return
	}

	// list
	fullQueryScope, err := query.MakeConditionFromQuery(apiQuery, t, query.GetAllConditionFilter())
	if err != nil {
		return nil, msg.NewError(reason.DatabaseError).WithError(err).WithStack()
	}
	var list []*T
	if err = txWithCtx.
		Scopes(crud.Preload()).
		Scopes(fullQueryScope).
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
