package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewUserRouter(
	handler *biz.UserHandler,
) *UserRouter {
	return &UserRouter{
		DefaultCrudRouter[*model.User]{
			BasePath: "/users",
			Handler:  handler,
		},
	}
}

type UserRouter struct {
	DefaultCrudRouter[*model.User]
}
