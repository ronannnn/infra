package repo

import (
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewJobGradeRepo() service.JobGradeRepo {
	return &jobGradeRepo{}
}

type jobGradeRepo struct {
	DefaultCrudRepo[model.JobGrade]
}
