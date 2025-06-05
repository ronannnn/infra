package biz

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewLoginRecordHandler(
	h handler.HttpHandler,
	srv *service.LoginRecordService,
) *LoginRecordHandler {
	return &LoginRecordHandler{
		handler.DefaultCrudHandler[*model.LoginRecord]{
			H:   h,
			Srv: srv,
		},
	}
}

type LoginRecordHandler struct {
	handler.DefaultCrudHandler[*model.LoginRecord]
}
