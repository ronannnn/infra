package repo

import (
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewLoginRecordRepo() service.LoginRecordRepo {
	return &loginRecordRepo{}
}

type loginRecordRepo struct {
	DefaultCrudRepo[*model.LoginRecord]
}
