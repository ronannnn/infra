package api

import (
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Service interface {
	Create(model *Api) error
	Update(partialUpdatedModel *Api) (Api, error)
	DeleteById(id uint) error
	DeleteByIds(ids []uint) error
	List(query ApiQuery) (response.PageResult, error)
	GetById(id uint) (Api, error)
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

func (srv *ServiceImpl) Create(model *Api) (err error) {
	return srv.store.create(srv.db, model)
}

func (srv *ServiceImpl) Update(partialUpdatedModel *Api) (Api, error) {
	return srv.store.update(srv.db, partialUpdatedModel)
}

func (srv *ServiceImpl) DeleteById(id uint) error {
	return srv.store.deleteById(srv.db, id)
}

func (srv *ServiceImpl) DeleteByIds(ids []uint) error {
	return srv.store.deleteByIds(srv.db, ids)
}

func (srv *ServiceImpl) List(query ApiQuery) (response.PageResult, error) {
	return srv.store.list(srv.db, query)
}

func (srv *ServiceImpl) GetById(id uint) (Api, error) {
	return srv.store.getById(srv.db, id)
}
