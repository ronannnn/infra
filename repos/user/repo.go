package user

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
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
		menuRepo: menuRepo,
		roleRepo: roleRepo,
	}
}

type repo struct {
	menuRepo menu.Repo
	roleRepo role.Repo
}

func (r repo) Create(tx *gorm.DB, model *models.User) error {
	return tx.Create(model).Error
}

func (r repo) Update(tx *gorm.DB, partialUpdatedModel *models.User) (updatedModel *models.User, err error) {
	if partialUpdatedModel.Id == 0 {
		return nil, msg.NewError(reason.DbModelUpdatedIdCannotBeZero).WithStack()
	}
	err = tx.Transaction(func(tx2 *gorm.DB) (err error) {
		if partialUpdatedModel.Roles != nil {
			if err = tx2.Model(partialUpdatedModel).Association("Roles").Unscoped().Replace(partialUpdatedModel.Roles); err != nil {
				return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
			}
			for _, item := range partialUpdatedModel.Roles {
				if _, err = r.roleRepo.Update(tx2, item); err != nil {
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
				if _, err = r.menuRepo.Update(tx2, item); err != nil {
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
		updatedModel, err = r.GetById(tx2, partialUpdatedModel.Id)
		return
	})
	return
}

func (r repo) DeleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.User{}, "id = ?", id).Error
}

func (r repo) DeleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&models.User{}, "id IN ?", ids).Error
}

func (r repo) List(tx *gorm.DB, userQuery query.Query) (result response.PageResult, err error) {
	var total int64
	var list []models.User
	if err = tx.Model(&models.User{}).Count(&total).Error; err != nil {
		return
	}
	queryScope, err := query.MakeConditionFromQuery(userQuery, models.User{})
	if err != nil {
		return
	}
	if err = tx.
		Scopes(queryScope).
		Scopes(query.Paginate(userQuery.Pagination)).
		Preload("Roles").
		Preload("Roles.Menus").
		Find(&list).Error; err != nil {
		return
	}
	result = response.PageResult{
		List:     list,
		Total:    total,
		PageNum:  userQuery.Pagination.PageNum,
		PageSize: userQuery.Pagination.PageSize,
	}
	return
}

func (r repo) GetById(tx *gorm.DB, id uint) (model *models.User, err error) {
	err = tx.Scopes(userPreload()).First(&model, "id = ?", id).Error
	return
}

func (r repo) GetByUsername(tx *gorm.DB, username string) (user *models.User, err error) {
	err = tx.Scopes(userPreload()).First(&user, &models.User{Username: &username}).Error
	return
}

func (r repo) GetByNickname(tx *gorm.DB, nickname string) (user *models.User, err error) {
	err = tx.Scopes(userPreload()).First(&user, &models.User{Nickname: &nickname}).Error
	return
}

func (r repo) ChangePwd(tx *gorm.DB, userId uint, newPwd string) error {
	return tx.Model(&models.User{}).Where("id = ?", userId).Update("password", newPwd).Error
}
