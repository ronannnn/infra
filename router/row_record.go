package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewRowRecordRouter(
	handler *biz.RowRecordHandler,
) *RowRecordRouter {
	return &RowRecordRouter{
		DefaultCrudRouter[*model.RowRecord]{
			Handler: handler,
		},
	}
}

type RowRecordRouter struct {
	DefaultCrudRouter[*model.RowRecord]
}

func (c *RowRecordRouter) GetBasePath() string {
	return "/row-records"
}
