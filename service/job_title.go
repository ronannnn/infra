package service

import (
	"github.com/ronannnn/infra/model"
	"gorm.io/gorm"
)

type JobTitleRepo interface {
	CrudRepo[model.JobTitle]
}

func NewJobTitleService(
	repo JobTitleRepo,
	db *gorm.DB,
) *JobTitleService {
	return &JobTitleService{
		DefaultCrudSrv[model.JobTitle]{
			Repo: repo,
			Db:   db,
		},
	}
}

type JobTitleService struct {
	DefaultCrudSrv[model.JobTitle]
}
