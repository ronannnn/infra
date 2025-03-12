package user

import (
	"context"

	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"github.com/ronannnn/infra/repos"
	"github.com/ronannnn/infra/services/menu"
	"github.com/ronannnn/infra/services/role"
	srv "github.com/ronannnn/infra/services/user"
	"gorm.io/gorm"
)

func New(
	menuRepo menu.Repo,
	roleRepo role.Repo,
) srv.Repo {
	return &repo{
		DefaultCrudRepo: repos.NewDefaultCrudRepo[models.User](
			"Roles",
			"Roles.Menus",
			"Menus",
		),
		menuRepo: menuRepo,
		roleRepo: roleRepo,
	}
}

type repo struct {
	repos.DefaultCrudRepo[models.User]
	menuRepo menu.Repo
	roleRepo role.Repo
}

func (r repo) Update(ctx context.Context, tx *gorm.DB, partialUpdatedModel *models.User) (updatedModel *models.User, err error) {
	if partialUpdatedModel.Id == 0 {
		return nil, msg.NewError(reason.DbModelUpdatedIdCannotBeZero).WithStack()
	}
	err = tx.Transaction(func(tx2 *gorm.DB) (err error) {
		if partialUpdatedModel.Roles != nil {
			if err = tx2.Model(partialUpdatedModel).Association("Roles").Unscoped().Replace(partialUpdatedModel.Roles); err != nil {
				return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
			}
			for _, item := range partialUpdatedModel.Roles {
				if _, err = r.roleRepo.Update(ctx, tx2, item); err != nil {
					return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
				}
			}
			partialUpdatedModel.Roles = nil
		}
		if partialUpdatedModel.Menus != nil {
			if err = tx2.Model(partialUpdatedModel).Association("Menus").Unscoped().Replace(partialUpdatedModel.Menus); err != nil {
				return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
			}
			for _, item := range partialUpdatedModel.Menus {
				if _, err = r.menuRepo.Update(ctx, tx2, item); err != nil {
					return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
				}
			}
			partialUpdatedModel.Menus = nil
		}
		result := tx2.Updates(partialUpdatedModel)
		if result.Error != nil {
			return msg.NewError(reason.DatabaseError).WithError(result.Error).WithStack()
		}
		if result.RowsAffected == 0 {
			return msg.NewError(reason.DbModelAlreadyUpdatedByOthers).WithStack()
		}
		updatedModel, err = r.GetById(ctx, tx2, partialUpdatedModel.Id)
		return
	})
	return
}

func (r repo) GetByUsername(ctx context.Context, tx *gorm.DB, username string) (user *models.User, err error) {
	err = tx.WithContext(ctx).Scopes(r.Preload()).First(&user, &models.User{Username: &username}).Error
	return
}

func (r repo) GetByNickname(ctx context.Context, tx *gorm.DB, nickname string) (user *models.User, err error) {
	err = tx.WithContext(ctx).Scopes(r.Preload()).First(&user, &models.User{Nickname: &nickname}).Error
	return
}

func (r repo) ChangePwd(ctx context.Context, tx *gorm.DB, userId uint, newPwd string) error {
	return tx.WithContext(ctx).Model(&models.User{}).Where("id = ?", userId).Update("password", newPwd).Error
}
