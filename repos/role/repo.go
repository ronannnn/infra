package role

import (
	"context"

	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"github.com/ronannnn/infra/repos"
	"github.com/ronannnn/infra/services/menu"
	srv "github.com/ronannnn/infra/services/role"
	"gorm.io/gorm"
)

func New(
	menuRepo menu.Repo,
) srv.Repo {
	return &repo{
		DefaultCrudRepo: repos.NewDefaultCrudRepo[*models.Role](),
		menuRepo:        menuRepo,
	}
}

type repo struct {
	repos.DefaultCrudRepo[*models.Role]
	menuRepo menu.Repo
}

func (r repo) Update(ctx context.Context, tx *gorm.DB, partialUpdatedModel *models.Role) (updatedModel *models.Role, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
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
			return models.ErrModified("models.Role")
		}
		return
	}); err != nil {
		return updatedModel, err
	}
	return r.GetById(ctx, tx, partialUpdatedModel.Id)
}

func (r repo) Preload() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Preload("Menus")
		return db
	}
}
