package biz

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewRowRecordHandler(
	h handler.HttpHandler,
	srv *service.RowRecordService,
) *RowRecordHandler {
	return &RowRecordHandler{
		handler.DefaultCrudHandler[model.RowRecord]{
			H:   h,
			Srv: srv,
		},
	}
}

type RowRecordHandler struct {
	handler.DefaultCrudHandler[model.RowRecord]
}
