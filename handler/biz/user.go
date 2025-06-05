package biz

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewUserHandler(
	h handler.HttpHandler,
	srv *service.UserService,
) *UserHandler {
	return &UserHandler{
		handler.DefaultCrudHandler[*model.User]{
			H:   h,
			Srv: srv,
		},
	}
}

type UserHandler struct {
	handler.DefaultCrudHandler[*model.User]
}
