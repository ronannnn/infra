package repo

import (
	"context"

	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"github.com/ronannnn/infra/service"
	"gorm.io/gorm"
)

func NewMenuRepo(
	apiRepo service.ApiRepo,
) service.MenuRepo {
	return &menuRepo{
		DefaultCrudRepo: NewDefaultCrudRepo[model.Menu]("Apis"),
		apiRepo:         apiRepo,
	}
}

type menuRepo struct {
	DefaultCrudRepo[model.Menu]
	apiRepo service.ApiRepo
}

func (r menuRepo) Update(ctx context.Context, tx *gorm.DB, partialUpdatedModel *model.Menu) (updatedModel *model.Menu, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, model.ErrUpdatedId
	}
	if err = tx.Transaction(func(tx *gorm.DB) (err error) {
		if partialUpdatedModel.Apis != nil {
			if err = tx.Model(partialUpdatedModel).Association("Apis").Unscoped().Replace(partialUpdatedModel.Apis); err != nil {
				return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
			}
			for _, item := range partialUpdatedModel.Apis {
				if _, err = r.apiRepo.Update(ctx, tx, item); err != nil {
					return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
				}
			}
			partialUpdatedModel.Apis = nil
		}
		result := tx.Updates(partialUpdatedModel)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return model.ErrModified("model.Menu")
		}
		return
	}); err != nil {
		return updatedModel, err
	}
	return r.GetById(ctx, tx, partialUpdatedModel.Id)
}
