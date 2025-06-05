package repo

import (
	"context"

	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"github.com/ronannnn/infra/service"
	"gorm.io/gorm"
)

func NewUserRepo(
	menuRepo service.MenuRepo,
	roleRepo service.RoleRepo,
) service.UserRepo {
	return &userRepo{
		DefaultCrudRepo: NewDefaultCrudRepo[*model.User](
			"Roles",
			"Roles.Menus",
			"Menus",
		),
		menuRepo: menuRepo,
		roleRepo: roleRepo,
	}
}

type userRepo struct {
	DefaultCrudRepo[*model.User]
	menuRepo service.MenuRepo
	roleRepo service.RoleRepo
}

func (r userRepo) Update(ctx context.Context, tx *gorm.DB, partialUpdatedModel *model.User) (updatedModel *model.User, err error) {
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

func (r userRepo) GetByUsername(ctx context.Context, tx *gorm.DB, username string) (user *model.User, err error) {
	err = tx.WithContext(ctx).Scopes(r.Preload()).First(&user, &model.User{Username: &username}).Error
	return
}

func (r userRepo) GetByNickname(ctx context.Context, tx *gorm.DB, nickname string) (user *model.User, err error) {
	err = tx.WithContext(ctx).Scopes(r.Preload()).First(&user, &model.User{Nickname: &nickname}).Error
	return
}

func (r userRepo) ChangePwd(ctx context.Context, tx *gorm.DB, userId uint, newPwd string) error {
	return tx.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId).Update("password", newPwd).Error
}
