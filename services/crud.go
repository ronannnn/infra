package services

import (
	"context"

	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type CrudRepo[T models.Crudable] interface {
	Create(context.Context, *gorm.DB, *T) error
	CreateWithScopes(context.Context, *gorm.DB, *T) (*T, error)
	Update(context.Context, *gorm.DB, *T) (*T, error)
	DeleteById(context.Context, *gorm.DB, uint) error
	DeleteByIds(context.Context, *gorm.DB, []uint) error
	List(context.Context, *gorm.DB, query.Query) (*response.PageResult, error)
	GetById(context.Context, *gorm.DB, uint) (*T, error)
}

type CrudSrv[T models.Crudable] interface {
	Create(context.Context, *T) error
	CreateWithScopes(context.Context, *T) (*T, error)
	Update(context.Context, *T) (*T, error)
	DeleteById(context.Context, uint) error
	DeleteByIds(context.Context, []uint) error
	List(context.Context, query.Query) (*response.PageResult, error)
	GetById(context.Context, uint) (*T, error)
}

type DefaultCrudSrv[T models.Crudable] struct {
	Db   *gorm.DB
	Repo CrudRepo[T]
}

func (srv *DefaultCrudSrv[T]) Create(ctx context.Context, model *T) error {
	return srv.Repo.Create(ctx, srv.Db, model)
}

func (srv *DefaultCrudSrv[T]) CreateWithScopes(ctx context.Context, model *T) (*T, error) {
	return srv.Repo.CreateWithScopes(ctx, srv.Db, model)
}

func (srv *DefaultCrudSrv[T]) Update(ctx context.Context, partialUpdatedModel *T) (*T, error) {
	return srv.Repo.Update(ctx, srv.Db, partialUpdatedModel)
}

func (srv *DefaultCrudSrv[T]) DeleteById(ctx context.Context, id uint) error {
	return srv.Repo.DeleteById(ctx, srv.Db, id)
}

func (srv *DefaultCrudSrv[T]) DeleteByIds(ctx context.Context, ids []uint) error {
	return srv.Repo.DeleteByIds(ctx, srv.Db, ids)
}

func (srv *DefaultCrudSrv[T]) List(ctx context.Context, query query.Query) (*response.PageResult, error) {
	return srv.Repo.List(ctx, srv.Db, query)
}

func (srv *DefaultCrudSrv[T]) GetById(ctx context.Context, id uint) (*T, error) {
	return srv.Repo.GetById(ctx, srv.Db, id)
}
