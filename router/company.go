package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewCompanyRouter(
	handler *biz.CompanyHandler,
) *CompanyRouter {
	return &CompanyRouter{
		DefaultCrudRouter[model.Company]{
			BasePath: "/companies",
			Handler:  handler,
		},
	}
}

type CompanyRouter struct {
	DefaultCrudRouter[model.Company]
}
