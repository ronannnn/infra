package wechattask

import (
	"context"

	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Service interface {
	Create(context.Context, *WechatTask) error
	Update(ctx context.Context, partialUpdatedModel *WechatTask) (WechatTask, error)
	DeleteById(ctx context.Context, id uint) error
	DeleteByIds(ctx context.Context, ids []uint) error
	List(ctx context.Context, query query.Query) (response.PageResult, error)
	GetById(ctx context.Context, id uint) (WechatTask, error)
	GetByUuid(ctx context.Context, uuid string) (WechatTask, error)
}

func ProvideService(
	db *gorm.DB,
	cfg *cfg.User,
	store Store,
) Service {
	return &ServiceImpl{
		db:    db,
		cfg:   cfg,
		store: store,
	}
}

type ServiceImpl struct {
	db    *gorm.DB
	cfg   *cfg.User
	store Store
}

func (srv *ServiceImpl) Create(ctx context.Context, model *WechatTask) (err error) {
	return srv.store.Create(srv.db, model)
}

func (srv *ServiceImpl) Update(ctx context.Context, partialUpdatedModel *WechatTask) (WechatTask, error) {
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

func (srv *ServiceImpl) GetById(ctx context.Context, id uint) (WechatTask, error) {
	return srv.store.GetById(srv.db, id)
}

func (srv *ServiceImpl) GetByUuid(ctx context.Context, uuid string) (WechatTask, error) {
	return srv.store.GetByUuid(srv.db, uuid)
}
