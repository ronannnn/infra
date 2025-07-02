package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewRowRecordRouter(
	handler *biz.RowRecordHandler,
) *RowRecordRouter {
	return &RowRecordRouter{
		DefaultCrudRouter[model.RowRecord]{
			BasePath: "/row-records",
			Handler:  handler,
		},
	}
}

type RowRecordRouter struct {
	DefaultCrudRouter[model.RowRecord]
}
