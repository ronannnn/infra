package loginrecord

import (
	"gorm.io/gorm"
)

type Store interface {
	Create(tx *gorm.DB, model *LoginRecord) error
}

func ProvideStore() Store {
	return StoreImpl{}
}

type StoreImpl struct {
}

func (s StoreImpl) Create(tx *gorm.DB, model *LoginRecord) error {
	return tx.Create(model).Error
}
