package biz

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewApiHandler(
	h handler.HttpHandler,
	srv *service.ApiService,
) *ApiHandler {
	return &ApiHandler{
		handler.DefaultCrudHandler[*model.Api]{
			H:   h,
			Srv: srv,
		},
	}
}

type ApiHandler struct {
	handler.DefaultCrudHandler[*model.Api]
}
