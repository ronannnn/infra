package company

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/api"
)

func ProvideHandler(
	h handler.HttpHandler,
	srv *api.Service,
) *Handler {
	return &Handler{
		handler.DefaultCrudHandler[*models.Api]{
			H:   h,
			Srv: srv,
		},
	}
}

type Handler struct {
	handler.DefaultCrudHandler[*models.Api]
}
