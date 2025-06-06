package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewApiRouter(
	handler *biz.ApiHandler,
) *ApiRouter {
	return &ApiRouter{
		DefaultCrudRouter[*model.Api]{
			Handler: handler,
		},
	}
}

type ApiRouter struct {
	DefaultCrudRouter[*model.Api]
}

func (c *ApiRouter) GetBasePath() string {
	return "/apis"
}
