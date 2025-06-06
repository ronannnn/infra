package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewDepartmentRouter(
	handler *biz.DepartmentHandler,
) *DepartmentRouter {
	return &DepartmentRouter{
		DefaultCrudRouter[*model.Department]{
			Handler: handler,
		},
	}
}

type DepartmentRouter struct {
	DefaultCrudRouter[*model.Department]
}

func (c *DepartmentRouter) GetBasePath() string {
	return "/departments"
}
