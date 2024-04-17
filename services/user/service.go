package user

import (
	"context"

	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Service interface {
	Create(context.Context, *models.User) error
	Update(ctx context.Context, partialUpdatedModel *models.User) (models.User, error)
	DeleteById(ctx context.Context, id uint) error
	DeleteByIds(ctx context.Context, ids []uint) error
	List(ctx context.Context, query query.UserQuery) (response.PageResult, error)
	GetById(ctx context.Context, id uint) (models.User, error)
	GetByNickname(ctx context.Context, nickname string) (models.User, error)
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

func (srv *ServiceImpl) Create(ctx context.Context, model *models.User) (err error) {
	if model.Password == nil {
		model.Password = &srv.cfg.DefaultHashedPassword
	}
	return srv.store.Create(srv.db, model)
}

func (srv *ServiceImpl) Update(ctx context.Context, partialUpdatedModel *models.User) (models.User, error) {
	return srv.store.Update(srv.db, partialUpdatedModel)
}

func (srv *ServiceImpl) DeleteById(ctx context.Context, id uint) error {
	return srv.store.DeleteById(srv.db, id)
}

func (srv *ServiceImpl) DeleteByIds(ctx context.Context, ids []uint) error {
	return srv.store.DeleteByIds(srv.db, ids)
}

func (srv *ServiceImpl) List(ctx context.Context, query query.UserQuery) (response.PageResult, error) {
	return srv.store.List(srv.db, query)
}

func (srv *ServiceImpl) GetById(ctx context.Context, id uint) (models.User, error) {
	return srv.store.GetById(srv.db, id)
}

func (srv *ServiceImpl) GetByNickname(ctx context.Context, nickname string) (models.User, error) {
	return srv.store.GetByNickname(srv.db, nickname)
}
