package menu

import (
	"context"

	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Service interface {
	Create(context.Context, *models.Menu) error
	Update(context.Context, *models.Menu) (models.Menu, error)
	DeleteById(context.Context, uint) error
	DeleteByIds(context.Context, []uint) error
	List(context.Context, query.Query) (response.PageResult, error)
	GetById(context.Context, uint) (models.Menu, error)
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

func (srv *ServiceImpl) Create(ctx context.Context, model *models.Menu) (err error) {
	return srv.store.Create(srv.db, model)
}

func (srv *ServiceImpl) Update(ctx context.Context, partialUpdatedModel *models.Menu) (models.Menu, error) {
	return srv.store.Update(srv.db, partialUpdatedModel)
}

func (srv *ServiceImpl) DeleteById(ctx context.Context, id uint) error {
	return srv.store.DeleteById(srv.db, id)
}

func (srv *ServiceImpl) DeleteByIds(ctx context.Context, ids []uint) error {
	return srv.store.DeleteByIds(srv.db, ids)
}

func (srv *ServiceImpl) List(ctx context.Context, query query.Query) (response.PageResult, error) {
	return srv.store.List(srv.db, query)
}

func (srv *ServiceImpl) GetById(ctx context.Context, id uint) (models.Menu, error) {
	return srv.store.GetById(srv.db, id)
}
