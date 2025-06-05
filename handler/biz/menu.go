package biz

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewMenuHandler(
	h handler.HttpHandler,
	srv *service.MenuService,
) *MenuHandler {
	return &MenuHandler{
		handler.DefaultCrudHandler[*model.Menu]{
			H:   h,
			Srv: srv,
		},
	}
}

type MenuHandler struct {
	handler.DefaultCrudHandler[*model.Menu]
}
