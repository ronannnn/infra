package biz

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewDepartmentHandler(
	h handler.HttpHandler,
	srv *service.DepartmentService,
) *DepartmentHandler {
	return &DepartmentHandler{
		handler.DefaultCrudHandler[*model.Department]{
			H:   h,
			Srv: srv,
		},
	}
}

type DepartmentHandler struct {
	handler.DefaultCrudHandler[*model.Department]
}
