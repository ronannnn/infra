package user

import (
	"context"

	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services"
	"gorm.io/gorm"
)

type Repo interface {
	services.CrudRepo[*models.User]
	GetByUsername(context.Context, *gorm.DB, string) (*models.User, error)
	GetByNickname(context.Context, *gorm.DB, string) (*models.User, error)
	ChangePwd(ctx context.Context, tx *gorm.DB, Id uint, newPwd string) error
}

func ProvideService(
	cfg *cfg.User,
	db *gorm.DB,
	repo Repo,
) *Service {
	return &Service{
		DefaultCrudSrv: services.DefaultCrudSrv[*models.User]{
			Db:   db,
			Repo: repo,
		},
		cfg:  cfg,
		repo: repo,
	}
}

type Service struct {
	services.DefaultCrudSrv[*models.User]
	cfg  *cfg.User
	repo Repo
}

func (srv *Service) Create(ctx context.Context, model *models.User) (err error) {
	if model.Password == nil {
		model.Password = &srv.cfg.DefaultHashedPassword
	}
	return srv.Repo.Create(ctx, srv.Db, model)
}

func (srv *Service) GetByNickname(ctx context.Context, nickname string) (*models.User, error) {
	return srv.repo.GetByNickname(ctx, srv.Db, nickname)
}

func (srv *Service) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return srv.repo.GetByUsername(ctx, srv.Db, username)
}
