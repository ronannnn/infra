package repo

import (
	"context"

	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"github.com/ronannnn/infra/service"
	"gorm.io/gorm"
)

func NewRoleRepo(
	menuRepo service.MenuRepo,
) service.RoleRepo {
	return &roleRepo{
		DefaultCrudRepo: NewDefaultCrudRepo[*model.Role]("Menus"),
		menuRepo:        menuRepo,
	}
}

type roleRepo struct {
	DefaultCrudRepo[*model.Role]
	menuRepo service.MenuRepo
}

func (r roleRepo) Update(ctx context.Context, tx *gorm.DB, partialUpdatedModel *model.Role) (updatedModel *model.Role, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, model.ErrUpdatedId
	}
	if err = tx.Transaction(func(tx *gorm.DB) (err error) {
		if partialUpdatedModel.Menus != nil {
			if err = tx.Model(partialUpdatedModel).Association("Menus").Unscoped().Replace(partialUpdatedModel.Menus); err != nil {
				return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
			}
			for _, item := range partialUpdatedModel.Menus {
				if _, err = r.menuRepo.Update(ctx, tx, item); err != nil {
					return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
				}
			}
			partialUpdatedModel.Menus = nil
		}
		result := tx.Updates(partialUpdatedModel)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return model.ErrModified("model.Role")
		}
		return
	}); err != nil {
		return updatedModel, err
	}
	return r.GetById(ctx, tx, partialUpdatedModel.Id)
}
