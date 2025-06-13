package repo

import (
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewApiRecordRepo() service.ApiRecordRepo {
	return &apiRecordRepo{}
}

type apiRecordRepo struct {
	DefaultCrudRepo[model.ApiRecord]
}
