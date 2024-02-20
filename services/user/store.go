package user

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Store interface {
	Create(model *models.User) error
	Update(partialUpdatedModel *models.User) (models.User, error)
	DeleteById(id uint) error
	DeleteByIds(ids []uint) error
	List(query query.UserQuery) (response.PageResult, error)
	GetById(id uint) (models.User, error)
	GetByUsername(username string) (models.User, error)
	ChangePwd(userId uint, newPwd string) error
}

func ProvideStore(db *gorm.DB) Store {
	return StoreImpl{db: db}
}

type StoreImpl struct {
	db *gorm.DB
}

func (s StoreImpl) Create(model *models.User) error {
	return s.db.Create(model).Error
}

func (s StoreImpl) Update(partialUpdatedModel *models.User) (updatedModel models.User, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	if err = s.db.Transaction(func(tx *gorm.DB) (err error) {
		// update associations with Associations()
		if partialUpdatedModel.Roles != nil {
			if err = tx.Model(partialUpdatedModel).Association("Roles").Replace(partialUpdatedModel.Roles); err != nil {
				return err
			}
			// set associations to nil to avoid Updates() below,
			partialUpdatedModel.Roles = nil
		}
		if partialUpdatedModel.Menus != nil {
			if err = tx.Model(partialUpdatedModel).Association("Menus").Replace(partialUpdatedModel.Menus); err != nil {
				return err
			}
			// set associations to nil to avoid Updates() below,
			partialUpdatedModel.Menus = nil
		}
		// update all other non-associations
		// if no other fields are updated, it still update the version so no error will occur
		result := tx.Updates(partialUpdatedModel)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return models.ErrModified("User")
		}
		return
	}); err != nil {
		return updatedModel, err
	}
	return s.GetById(partialUpdatedModel.Id)
}

func (s StoreImpl) DeleteById(id uint) error {
	return s.db.Delete(&models.User{}, "id = ?", id).Error
}

func (s StoreImpl) DeleteByIds(ids []uint) error {
	return s.db.Delete(&models.User{}, "id IN ?", ids).Error
}

func (s StoreImpl) List(userQuery query.UserQuery) (result response.PageResult, err error) {
	var total int64
	var list []models.User
	if err = s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return
	}
	if err = s.db.
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

func (s StoreImpl) GetById(id uint) (model models.User, err error) {
	err = s.db.Preload("Roles").Preload("Roles.Menus").Preload("Menus").First(&model, "id = ?", id).Error
	return
}

func (s StoreImpl) GetByUsername(username string) (user models.User, err error) {
	err = s.db.First(&user, &models.User{Username: &username}).Error
	return
}

func (s StoreImpl) ChangePwd(userId uint, newPwd string) error {
	return s.db.Model(&models.User{}).Where("id = ?", userId).Update("password", newPwd).Error
}
