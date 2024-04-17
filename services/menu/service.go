package menu

import (
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Service interface {
	Create(model *Menu) error
	Update(partialUpdatedModel *Menu) (Menu, error)
	DeleteById(id uint) error
	DeleteByIds(ids []uint) error
	List(query MenuQuery) (response.PageResult, error)
	GetById(id uint) (Menu, error)
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

func (srv *ServiceImpl) Create(model *Menu) (err error) {
	return srv.store.Create(srv.db, model)
}

func (srv *ServiceImpl) Update(partialUpdatedModel *Menu) (Menu, error) {
	return srv.store.Update(srv.db, partialUpdatedModel)
}

func (srv *ServiceImpl) DeleteById(id uint) error {
	return srv.store.DeleteById(srv.db, id)
}

func (srv *ServiceImpl) DeleteByIds(ids []uint) error {
	return srv.store.DeleteByIds(srv.db, ids)
}

func (srv *ServiceImpl) List(query MenuQuery) (response.PageResult, error) {
	return srv.store.List(srv.db, query)
}

func (srv *ServiceImpl) GetById(id uint) (Menu, error) {
	return srv.store.GetById(srv.db, id)
}
