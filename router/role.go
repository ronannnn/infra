package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewRoleRouter(
	handler *biz.RoleHandler,
) *RoleRouter {
	return &RoleRouter{
		DefaultCrudRouter[*model.Role]{
			BasePath: "/roles",
			Handler:  handler,
		},
	}
}

type RoleRouter struct {
	DefaultCrudRouter[*model.Role]
}
