package biz

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewRoleHandler(
	h handler.HttpHandler,
	srv *service.RoleService,
) *RoleHandler {
	return &RoleHandler{
		handler.DefaultCrudHandler[model.Role]{
			H:   h,
			Srv: srv,
		},
	}
}

type RoleHandler struct {
	handler.DefaultCrudHandler[model.Role]
}
