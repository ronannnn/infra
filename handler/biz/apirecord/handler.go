package biz

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/apirecord"
)

func ProvideHandler(
	h handler.HttpHandler,
	srv *apirecord.Service,
) *Handler {
	return &Handler{
		handler.DefaultCrudHandler[*models.ApiRecord]{
			H:   h,
			Srv: srv,
		},
	}
}

type Handler struct {
	handler.DefaultCrudHandler[*models.ApiRecord]
}
