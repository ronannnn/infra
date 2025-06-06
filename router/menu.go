package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewMenuRouter(
	handler *biz.MenuHandler,
) *MenuRouter {
	return &MenuRouter{
		DefaultCrudRouter[*model.Menu]{
			Handler: handler,
		},
	}
}

type MenuRouter struct {
	DefaultCrudRouter[*model.Menu]
}

func (c *MenuRouter) GetBasePath() string {
	return "/menus"
}
