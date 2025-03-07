package menu

import (
	"context"

	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Repo interface {
	Create(*gorm.DB, *models.Menu) error
	Update(*gorm.DB, *models.Menu) (models.Menu, error)
	DeleteById(*gorm.DB, uint) error
	DeleteByIds(*gorm.DB, []uint) error
	List(*gorm.DB, query.Query) (response.PageResult, error)
	GetById(*gorm.DB, uint) (models.Menu, error)
}

func ProvideService(
	repo Repo,
	db *gorm.DB,
) *Service {
	return &Service{
		repo: repo,
		db:   db,
	}
}

type Service struct {
	repo Repo
	db   *gorm.DB
}

func (srv *Service) Create(ctx context.Context, model *models.Menu) (err error) {
	return srv.repo.Create(srv.db.WithContext(ctx), model)
}

func (srv *Service) Update(ctx context.Context, partialUpdatedModel *models.Menu) (models.Menu, error) {
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

func (srv *Service) GetById(ctx context.Context, id uint) (models.Menu, error) {
	return srv.repo.GetById(srv.db.WithContext(ctx), id)
}
