package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewDepartmentRouter(
	handler *biz.DepartmentHandler,
) *DepartmentRouter {
	return &DepartmentRouter{
		DefaultCrudRouter[model.Department]{
			BasePath: "/departments",
			Handler:  handler,
		},
	}
}

type DepartmentRouter struct {
	DefaultCrudRouter[model.Department]
}
