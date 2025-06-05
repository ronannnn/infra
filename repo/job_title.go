package repo

import (
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewJobTitleRepo() service.JobTitleRepo {
	return &jobTitleRepo{}
}

type jobTitleRepo struct {
	DefaultCrudRepo[*model.JobTitle]
}
