package biz

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewCompanyHandler(
	h handler.HttpHandler,
	srv *service.CompanyService,
) *CompanyHandler {
	return &CompanyHandler{
		handler.DefaultCrudHandler[*model.Company]{
			H:   h,
			Srv: srv,
		},
	}
}

type CompanyHandler struct {
	handler.DefaultCrudHandler[*model.Company]
}
