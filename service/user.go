package service

import (
	"context"

	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/model/request/query"
	"github.com/ronannnn/infra/model/response"
	"gorm.io/gorm"
)

type UserRepo interface {
	Create(context.Context, *gorm.DB, *model.User) error
	CreateWithScopes(context.Context, *gorm.DB, *model.User) (*model.User, error)
	Update(context.Context, *gorm.DB, *model.User) (*model.User, error)
	DeleteById(context.Context, *gorm.DB, uint) error
	DeleteByIds(context.Context, *gorm.DB, []uint) error
	List(context.Context, *gorm.DB, query.Query) (*response.PageResult, error)
	GetById(context.Context, *gorm.DB, uint) (*model.User, error)

	GetByUsername(context.Context, *gorm.DB, string) (*model.User, error)
	GetByNickname(context.Context, *gorm.DB, string) (*model.User, error)
	ChangePwd(ctx context.Context, tx *gorm.DB, Id uint, newPwd string) error
}

func NewUserService(
	cfg *cfg.User,
	db *gorm.DB,
	repo UserRepo,
) *UserService {
	return &UserService{
		cfg:  cfg,
		db:   db,
		repo: repo,
	}
}

type UserService struct {
	cfg  *cfg.User
	db   *gorm.DB
	repo UserRepo
}

func (srv *UserService) Create(ctx context.Context, model *model.User) (err error) {
	if model.Password == nil {
		model.Password = &srv.cfg.DefaultHashedPassword
	}
	return srv.repo.Create(ctx, srv.db, model)
}

func (srv *UserService) CreateWithScopes(ctx context.Context, model *model.User) (*model.User, error) {
	return srv.repo.CreateWithScopes(ctx, srv.db, model)
}

func (srv *UserService) Update(ctx context.Context, model *model.User) (*model.User, error) {
	return srv.repo.Update(ctx, srv.db, model)
}

func (srv *UserService) DeleteById(ctx context.Context, id uint) error {
	return srv.repo.DeleteById(ctx, srv.db, id)
}

func (srv *UserService) DeleteByIds(ctx context.Context, ids []uint) error {
	return srv.repo.DeleteByIds(ctx, srv.db, ids)
}

func (srv *UserService) List(ctx context.Context, q query.Query) (*response.PageResult, error) {
	return srv.repo.List(ctx, srv.db, q)
}

func (srv *UserService) GetById(ctx context.Context, id uint) (*model.User, error) {
	return srv.repo.GetById(ctx, srv.db, id)
}

func (srv *UserService) GetByNickname(ctx context.Context, nickname string) (*model.User, error) {
	return srv.repo.GetByNickname(ctx, srv.db, nickname)
}

func (srv *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return srv.repo.GetByUsername(ctx, srv.db, username)
}
