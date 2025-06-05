package biz

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/company"
)

func ProvideHandler(
	h handler.HttpHandler,
	srv *company.Service,
) *Handler {
	return &Handler{
		handler.DefaultCrudHandler[*models.Company]{
			H:   h,
			Srv: srv,
		},
	}
}

type Handler struct {
	handler.DefaultCrudHandler[*models.Company]
}
