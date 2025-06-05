package service

import (
	"context"

	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/model"
	"gorm.io/gorm"
)

type UserRepo interface {
	CrudRepo[*model.User]
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
		DefaultCrudSrv: DefaultCrudSrv[*model.User]{
			Db:   db,
			Repo: repo,
		},
		cfg:  cfg,
		repo: repo,
	}
}

type UserService struct {
	DefaultCrudSrv[*model.User]
	cfg  *cfg.User
	repo UserRepo
}

func (srv *UserService) Create(ctx context.Context, model *model.User) (err error) {
	if model.Password == nil {
		model.Password = &srv.cfg.DefaultHashedPassword
	}
	return srv.Repo.Create(ctx, srv.Db, model)
}

func (srv *UserService) GetByNickname(ctx context.Context, nickname string) (*model.User, error) {
	return srv.repo.GetByNickname(ctx, srv.Db, nickname)
}

func (srv *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return srv.repo.GetByUsername(ctx, srv.Db, username)
}
