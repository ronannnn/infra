package biz

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewApiRecordHandler(
	h handler.HttpHandler,
	srv *service.ApiRecordService,
) *ApiRecordHandler {
	return &ApiRecordHandler{
		handler.DefaultCrudHandler[*model.ApiRecord]{
			H:   h,
			Srv: srv,
		},
	}
}

type ApiRecordHandler struct {
	handler.DefaultCrudHandler[*model.ApiRecord]
}
