package apirecord

import (
	"context"

	"github.com/ronannnn/infra/models"
	"gorm.io/gorm"
)

type Repo interface {
	Create(*gorm.DB, *models.ApiRecord) error
	Update(*gorm.DB, *models.ApiRecord) (*models.ApiRecord, error)
	DeleteById(*gorm.DB, uint) error
	DeleteByIds(*gorm.DB, []uint) error
	GetById(*gorm.DB, uint) (*models.ApiRecord, error)
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

func (srv *Service) Create(ctx context.Context, model *models.ApiRecord) error {
	return srv.repo.Create(srv.db.WithContext(ctx), model)
}

func (srv *Service) Update(ctx context.Context, partialUpdatedModel *models.ApiRecord) (updatedModel *models.ApiRecord, err error) {
	return srv.repo.Update(srv.db.WithContext(ctx), partialUpdatedModel)
}

func (srv *Service) DeleteById(ctx context.Context, id uint) error {
	return srv.repo.DeleteById(srv.db.WithContext(ctx), id)
}

func (srv *Service) DeleteByIds(ctx context.Context, ids []uint) error {
	return srv.repo.DeleteByIds(srv.db.WithContext(ctx), ids)
}
