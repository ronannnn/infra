package user

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	srv "github.com/ronannnn/infra/services/user"
	"gorm.io/gorm"
)

func New() srv.Repo {
	return &repo{}
}

type repo struct {
}

func (s repo) Create(tx *gorm.DB, model *models.User) error {
	return tx.Create(model).Error
}

func (s repo) Update(tx *gorm.DB, partialUpdatedModel *models.User) (updatedModel models.User, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	if err = tx.Updates(partialUpdatedModel).Error; err != nil {
		return updatedModel, err
	}
	return s.GetById(tx, partialUpdatedModel.Id)
}

func (s repo) DeleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.User{}, "id = ?", id).Error
}

func (s repo) DeleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&models.User{}, "id IN ?", ids).Error
}

func (s repo) List(tx *gorm.DB, userQuery query.Query) (result response.PageResult, err error) {
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

func (s repo) GetById(tx *gorm.DB, id uint) (model models.User, err error) {
	err = tx.Scopes(userPreload()).First(&model, "id = ?", id).Error
	return
}

func (s repo) GetByUsername(tx *gorm.DB, username string) (user models.User, err error) {
	err = tx.Scopes(userPreload()).First(&user, &models.User{Username: &username}).Error
	return
}

func (s repo) GetByNickname(tx *gorm.DB, nickname string) (user models.User, err error) {
	err = tx.Scopes(userPreload()).First(&user, &models.User{Nickname: &nickname}).Error
	return
}

func (s repo) ChangePwd(tx *gorm.DB, userId uint, newPwd string) error {
	return tx.Model(&models.User{}).Where("id = ?", userId).Update("password", newPwd).Error
}
