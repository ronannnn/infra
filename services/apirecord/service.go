package apirecord

import "gorm.io/gorm"

type Service interface {
	Create(model *ApiRecord) error
	Update(partialUpdatedModel *ApiRecord) (ApiRecord, error)
	DeleteById(id uint) error
	DeleteByIds(ids []uint) error
}

func ProvideService(
	store Store,
	db *gorm.DB,
) Service {
	return &ServiceImpl{
		store: store,
		db:    db,
	}
}

type ServiceImpl struct {
	store Store
	db    *gorm.DB
}

func (srv *ServiceImpl) Create(model *ApiRecord) error {
	return srv.store.create(srv.db, model)
}

func (srv *ServiceImpl) Update(partialUpdatedModel *ApiRecord) (updatedModel ApiRecord, err error) {
	return srv.store.update(srv.db, partialUpdatedModel)
}

func (srv *ServiceImpl) DeleteById(id uint) error {
	return srv.store.deleteById(srv.db, id)
}

func (srv *ServiceImpl) DeleteByIds(ids []uint) error {
	return srv.store.deleteByIds(srv.db, ids)
}
