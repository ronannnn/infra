package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewMenuRouter(
	handler *biz.MenuHandler,
) *MenuRouter {
	return &MenuRouter{
		DefaultCrudRouter[model.Menu]{
			BasePath: "/menus",
			Handler:  handler,
		},
	}
}

type MenuRouter struct {
	DefaultCrudRouter[model.Menu]
}
