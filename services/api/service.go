package api

import (
	"context"

	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Service interface {
	Create(context.Context, *models.Api) error
	Update(context.Context, *models.Api) (models.Api, error)
	DeleteById(context.Context, uint) error
	DeleteByIds(context.Context, []uint) error
	List(context.Context, query.Query) (response.PageResult, error)
	GetById(context.Context, uint) (models.Api, error)
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

func (srv *ServiceImpl) Create(ctx context.Context, model *models.Api) (err error) {
	return srv.store.create(srv.db, model)
}

func (srv *ServiceImpl) Update(ctx context.Context, partialUpdatedModel *models.Api) (models.Api, error) {
	return srv.store.update(srv.db, partialUpdatedModel)
}

func (srv *ServiceImpl) DeleteById(ctx context.Context, id uint) error {
	return srv.store.deleteById(srv.db, id)
}

func (srv *ServiceImpl) DeleteByIds(ctx context.Context, ids []uint) error {
	return srv.store.deleteByIds(srv.db, ids)
}

func (srv *ServiceImpl) List(ctx context.Context, query query.Query) (response.PageResult, error) {
	return srv.store.list(srv.db, query)
}

func (srv *ServiceImpl) GetById(ctx context.Context, id uint) (models.Api, error) {
	return srv.store.getById(srv.db, id)
}
