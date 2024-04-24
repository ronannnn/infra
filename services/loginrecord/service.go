package loginrecord

import (
	"gorm.io/gorm"
)

type Service interface {
	Create(model *LoginRecord) error
}

func ProvideService(
	db *gorm.DB,
	store Store,
) Service {
	return &ServiceImpl{
		db:    db,
		store: store,
	}
}

type ServiceImpl struct {
	db    *gorm.DB
	store Store
}

func (srv *ServiceImpl) Create(model *LoginRecord) (err error) {
	return srv.store.Create(srv.db, model)
}
