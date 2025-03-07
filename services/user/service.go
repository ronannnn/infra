package user

import (
	"context"

	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Repo interface {
	Create(*gorm.DB, *models.User) error
	Update(*gorm.DB, *models.User) (models.User, error)
	DeleteById(*gorm.DB, uint) error
	DeleteByIds(*gorm.DB, []uint) error
	List(*gorm.DB, query.Query) (response.PageResult, error)
	GetById(*gorm.DB, uint) (models.User, error)
	GetByUsername(*gorm.DB, string) (models.User, error)
	GetByNickname(*gorm.DB, string) (models.User, error)
	ChangePwd(tx *gorm.DB, Id uint, newPwd string) error
}

func ProvideService(
	cfg *cfg.User,
	db *gorm.DB,
	repo Repo,
) *Service {
	return &Service{
		db:   db,
		cfg:  cfg,
		repo: repo,
	}
}

type Service struct {
	cfg  *cfg.User
	db   *gorm.DB
	repo Repo
}

func (srv *Service) Create(ctx context.Context, model *models.User) (err error) {
	if model.Password == nil {
		model.Password = &srv.cfg.DefaultHashedPassword
	}
	return srv.repo.Create(srv.db.WithContext(ctx), model)
}

func (srv *Service) Update(ctx context.Context, partialUpdatedModel *models.User) (models.User, error) {
	return srv.repo.Update(srv.db.WithContext(ctx), partialUpdatedModel)
}

func (srv *Service) DeleteById(ctx context.Context, id uint) error {
	return srv.repo.DeleteById(srv.db.WithContext(ctx), id)
}

func (srv *Service) DeleteByIds(ctx context.Context, ids []uint) error {
	return srv.repo.DeleteByIds(srv.db.WithContext(ctx), ids)
}

func (srv *Service) List(ctx context.Context, query query.Query) (response.PageResult, error) {
	return srv.repo.List(srv.db.WithContext(ctx), query)
}

func (srv *Service) GetById(ctx context.Context, id uint) (models.User, error) {
	return srv.repo.GetById(srv.db.WithContext(ctx), id)
}

func (srv *Service) GetByNickname(ctx context.Context, nickname string) (models.User, error) {
	return srv.repo.GetByNickname(srv.db.WithContext(ctx), nickname)
}

func (srv *Service) GetByUsername(ctx context.Context, username string) (models.User, error) {
	return srv.repo.GetByUsername(srv.db.WithContext(ctx), username)
}
