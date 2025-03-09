package menu

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services"
	"gorm.io/gorm"
)

type Repo interface {
	services.CrudRepo[models.Menu]
}

func ProvideService(
	repo Repo,
	db *gorm.DB,
) *Service {
	return &Service{
		services.DefaultCrudSrv[models.Menu]{
			Repo: repo,
			Db:   db,
		},
	}
}

type Service struct {
	services.DefaultCrudSrv[models.Menu]
}
