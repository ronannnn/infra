package useropenid

import (
	"fmt"

	"gorm.io/gorm"
)

type Store interface {
	create(*gorm.DB, *UserOpenId) error
	updateOpenIdByUserIdAndAppId(*gorm.DB, *UserOpenId) error
	getByUserIdAndAppId(*gorm.DB, uint, string) (UserOpenId, error)
}

func ProvideStore() Store {
	return StoreImpl{}
}

type StoreImpl struct {
}

func (s StoreImpl) create(tx *gorm.DB, model *UserOpenId) error {
	return tx.Create(model).Error
}

func (s StoreImpl) updateOpenIdByUserIdAndAppId(tx *gorm.DB, newModel *UserOpenId) (err error) {
	if newModel.UserId == 0 || newModel.OfficialAccountAppId == "" {
		return fmt.Errorf("UserId and OfficialAccountAppId are required")
	}
	err = tx.
		Where(&UserOpenId{UserId: newModel.UserId, OfficialAccountAppId: newModel.OfficialAccountAppId}).
		Updates(UserOpenId{OpenId: newModel.OpenId}).Error
	return
}

func (s StoreImpl) getByUserIdAndAppId(tx *gorm.DB, userId uint, appId string) (model UserOpenId, err error) {
	err = tx.
		Where(&UserOpenId{UserId: userId, OfficialAccountAppId: appId}).
		First(&model).
		Error
	return
}
