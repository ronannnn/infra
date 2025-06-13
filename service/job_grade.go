package service

import (
	"github.com/ronannnn/infra/model"
	"gorm.io/gorm"
)

type JobGradeRepo interface {
	CrudRepo[model.JobGrade]
}

func NewJobGradeService(
	repo JobGradeRepo,
	db *gorm.DB,
) *JobGradeService {
	return &JobGradeService{
		DefaultCrudSrv[model.JobGrade]{
			Repo: repo,
			Db:   db,
		},
	}
}

type JobGradeService struct {
	DefaultCrudSrv[model.JobGrade]
}
