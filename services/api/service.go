package api

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Service interface {
	Create(model *models.Api) error
	Update(partialUpdatedModel *models.Api) (models.Api, error)
	DeleteById(id uint) error
	DeleteByIds(ids []uint) error
	List(query query.Query) (response.PageResult, error)
	GetById(id uint) (models.Api, error)
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

func (srv *ServiceImpl) Create(model *models.Api) (err error) {
	return srv.store.create(srv.db, model)
}

func (srv *ServiceImpl) Update(partialUpdatedModel *models.Api) (models.Api, error) {
	return srv.store.update(srv.db, partialUpdatedModel)
}

func (srv *ServiceImpl) DeleteById(id uint) error {
	return srv.store.deleteById(srv.db, id)
}

func (srv *ServiceImpl) DeleteByIds(ids []uint) error {
	return srv.store.deleteByIds(srv.db, ids)
}

func (srv *ServiceImpl) List(query query.Query) (response.PageResult, error) {
	return srv.store.list(srv.db, query)
}

func (srv *ServiceImpl) GetById(id uint) (models.Api, error) {
	return srv.store.getById(srv.db, id)
}
