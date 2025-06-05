package repo

import (
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewApiRepo() service.ApiRepo {
	return &apiRepo{}
}

type apiRepo struct {
	DefaultCrudRepo[*model.Api]
}
