package user

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Store interface {
	Create(tx *gorm.DB, model *models.User) error
	Update(tx *gorm.DB, partialUpdatedModel *models.User) (models.User, error)
	DeleteById(tx *gorm.DB, id uint) error
	DeleteByIds(tx *gorm.DB, ids []uint) error
	List(tx *gorm.DB, query query.UserQuery) (response.PageResult, error)
	GetById(tx *gorm.DB, id uint) (models.User, error)
	GetByUsername(tx *gorm.DB, username string) (models.User, error)
	GetByNickname(tx *gorm.DB, nickname string) (models.User, error)
	ChangePwd(tx *gorm.DB, userId uint, newPwd string) error
}

func ProvideStore() Store {
	return StoreImpl{}
}

type StoreImpl struct {
}

func (s StoreImpl) Create(tx *gorm.DB, model *models.User) error {
	return tx.Create(model).Error
}

func (s StoreImpl) Update(tx *gorm.DB, partialUpdatedModel *models.User) (updatedModel models.User, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	if err = tx.Updates(partialUpdatedModel).Error; err != nil {
		return updatedModel, err
	}
	return s.GetById(tx, partialUpdatedModel.Id)
}

func (s StoreImpl) DeleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.User{}, "id = ?", id).Error
}

func (s StoreImpl) DeleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&models.User{}, "id IN ?", ids).Error
}

func (s StoreImpl) List(tx *gorm.DB, userQuery query.UserQuery) (result response.PageResult, err error) {
	var total int64
	var list []models.User
	if err = tx.Model(&models.User{}).Count(&total).Error; err != nil {
		return
	}
	if err = tx.
		Scopes(query.MakeConditionFromQuery(userQuery)).
		Scopes(query.Paginate(userQuery.Pagination.PageNum, userQuery.Pagination.PageSize)).
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

func (s StoreImpl) GetById(tx *gorm.DB, id uint) (model models.User, err error) {
	err = tx.Preload("Roles").Preload("Roles.Menus").Preload("Menus").First(&model, "id = ?", id).Error
	return
}

func (s StoreImpl) GetByUsername(tx *gorm.DB, username string) (user models.User, err error) {
	err = tx.First(&user, &models.User{Username: &username}).Error
	return
}

func (s StoreImpl) GetByNickname(tx *gorm.DB, nickname string) (user models.User, err error) {
	err = tx.First(&user, &models.User{Nickname: &nickname}).Error
	return
}

func (s StoreImpl) ChangePwd(tx *gorm.DB, userId uint, newPwd string) error {
	return tx.Model(&models.User{}).Where("id = ?", userId).Update("password", newPwd).Error
}
