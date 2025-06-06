package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewApiRecordRouter(
	handler *biz.ApiHandler,
) *ApiRecordRouter {
	return &ApiRecordRouter{
		DefaultCrudRouter[*model.ApiRecord]{
			Handler: handler,
		},
	}
}

type ApiRecordRouter struct {
	DefaultCrudRouter[*model.ApiRecord]
}

func (c *ApiRecordRouter) GetBasePath() string {
	return "/api-records"
}
