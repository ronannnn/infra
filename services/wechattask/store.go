package wechattask

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Store interface {
	Create(tx *gorm.DB, model *WechatTask) error
	Update(tx *gorm.DB, partialUpdatedModel *WechatTask) (WechatTask, error)
	DeleteById(tx *gorm.DB, id uint) error
	DeleteByIds(tx *gorm.DB, ids []uint) error
	List(tx *gorm.DB, query WechatTaskQuery) (response.PageResult, error)
	GetById(tx *gorm.DB, id uint) (WechatTask, error)
	GetByUuid(tx *gorm.DB, uuid string) (WechatTask, error)
}

func ProvideStore(
	wechatTaskUserIdStore WechatTaskUserIdStore,
) Store {
	return StoreImpl{
		wechatTaskUserIdStore: wechatTaskUserIdStore,
	}
}

type StoreImpl struct {
	wechatTaskUserIdStore WechatTaskUserIdStore
}

func (s StoreImpl) Create(tx *gorm.DB, model *WechatTask) error {
	return tx.Create(model).Error
}

func (s StoreImpl) Update(tx *gorm.DB, partialUpdatedModel *WechatTask) (updatedModel WechatTask, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	err = tx.Transaction(func(tx2 *gorm.DB) (err error) {
		if partialUpdatedModel.WechatTaskUserIds != nil {
			if err = tx2.Model(partialUpdatedModel).Association("WechatTaskUserIds").Unscoped().Replace(partialUpdatedModel.WechatTaskUserIds); err != nil {
				return
			}
			for _, WechatTaskUserId := range *partialUpdatedModel.WechatTaskUserIds {
				if _, err = s.wechatTaskUserIdStore.update(tx2, &WechatTaskUserId); err != nil {
					return
				}
			}
			partialUpdatedModel.WechatTaskUserIds = nil
		}
		result := tx2.Updates(partialUpdatedModel)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return models.ErrModified("Wechat Task")
		}
		updatedModel, err = s.GetById(tx2, partialUpdatedModel.Id)
		return
	})
	return
}

func (s StoreImpl) DeleteById(tx *gorm.DB, id uint) error {
	return tx.Delete(&WechatTask{}, "id = ?", id).Error
}

func (s StoreImpl) DeleteByIds(tx *gorm.DB, ids []uint) error {
	return tx.Delete(&WechatTask{}, "id IN ?", ids).Error
}

func (s StoreImpl) List(tx *gorm.DB, userQuery WechatTaskQuery) (result response.PageResult, err error) {
	var total int64
	var list []WechatTask
	if err = tx.Model(&WechatTask{}).Count(&total).Error; err != nil {
		return
	}
	if err = tx.
		Scopes(query.MakeConditionFromQuery(userQuery)).
		Scopes(query.Paginate(userQuery.Pagination.PageNum, userQuery.Pagination.PageSize)).
		Scopes(wechatTaskPreload()).
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

func (s StoreImpl) GetById(tx *gorm.DB, id uint) (model WechatTask, err error) {
	err = tx.Scopes(wechatTaskPreload()).First(&model, "id = ?", id).Error
	return
}

func (s StoreImpl) GetByUuid(tx *gorm.DB, uuid string) (model WechatTask, err error) {
	err = tx.Scopes(wechatTaskPreload()).First(&model, "uuid = ?", uuid).Error
	return
}

type WechatTaskUserIdStore interface {
	update(tx *gorm.DB, partialUpdatedModel *WechatTaskUserId) (WechatTaskUserId, error)
	getById(tx *gorm.DB, id uint) (WechatTaskUserId, error)
}

func ProvideWechatTaskUserIdStore() WechatTaskUserIdStore {
	return &WechatTaskUserIdImpl{}
}

type WechatTaskUserIdImpl struct {
}

func (s *WechatTaskUserIdImpl) update(tx *gorm.DB, partialUpdatedModel *WechatTaskUserId) (updatedModel WechatTaskUserId, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	err = tx.Transaction(func(tx2 *gorm.DB) (err error) {
		result := tx2.Updates(partialUpdatedModel)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return models.ErrModified("Wechat Task User Id")
		}
		updatedModel, err = s.getById(tx2, partialUpdatedModel.Id)
		return
	})
	return
}

func (s *WechatTaskUserIdImpl) getById(tx *gorm.DB, id uint) (model WechatTaskUserId, err error) {
	err = tx.First(&model, "id = ?", id).Error
	return
}
