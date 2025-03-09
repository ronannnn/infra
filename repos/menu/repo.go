package menu

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"github.com/ronannnn/infra/repos"
	"github.com/ronannnn/infra/services/api"
	srv "github.com/ronannnn/infra/services/menu"
	"gorm.io/gorm"
)

func New(
	apiRepo api.Repo,
) srv.Repo {
	return &repo{
		DefaultCrudRepo: repos.NewDefaultCrudRepo[models.Menu](),
		apiRepo:         apiRepo,
	}
}

type repo struct {
	repos.DefaultCrudRepo[models.Menu]
	apiRepo api.Repo
}

func (r repo) Update(tx *gorm.DB, partialUpdatedModel *models.Menu) (updatedModel *models.Menu, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	if err = tx.Transaction(func(tx *gorm.DB) (err error) {
		if partialUpdatedModel.Apis != nil {
			if err = tx.Model(partialUpdatedModel).Association("Apis").Unscoped().Replace(partialUpdatedModel.Apis); err != nil {
				return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
			}
			for _, item := range partialUpdatedModel.Apis {
				if _, err = r.apiRepo.Update(tx, item); err != nil {
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
			return models.ErrModified("models.Menu")
		}
		return
	}); err != nil {
		return updatedModel, err
	}
	return r.GetById(tx, partialUpdatedModel.Id)
}

func (r repo) Preload() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Preload("Apis")
		return db
	}
}
