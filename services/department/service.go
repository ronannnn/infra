package department

import (
	"context"

	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Repo interface {
	Create(*gorm.DB, *models.Department) error
	Update(*gorm.DB, *models.Department) (*models.Department, error)
	DeleteById(*gorm.DB, uint) error
	DeleteByIds(*gorm.DB, []uint) error
	List(*gorm.DB, query.Query) (response.PageResult, error)
	GetById(*gorm.DB, uint) (*models.Department, error)
}

func ProvideService(
	db *gorm.DB,
	repo Repo,
) *Service {
	return &Service{
		db:   db,
		repo: repo,
	}
}

type Service struct {
	db   *gorm.DB
	repo Repo
}

func (srv *Service) Create(ctx context.Context, model *models.Department) (err error) {
	return srv.repo.Create(srv.db.WithContext(ctx), model)
}

func (srv *Service) Update(ctx context.Context, partialUpdatedModel *models.Department) (*models.Department, error) {
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

func (srv *Service) GetById(ctx context.Context, id uint) (*models.Department, error) {
	return srv.repo.GetById(srv.db.WithContext(ctx), id)
}
