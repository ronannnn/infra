package loginrecord

import (
	"github.com/ronannnn/infra/models"
	"gorm.io/gorm"
)

type Repo interface {
	Create(*gorm.DB, *models.LoginRecord) error
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

func (srv *Service) Create(model *models.LoginRecord) (err error) {
	return srv.repo.Create(srv.db, model)
}
