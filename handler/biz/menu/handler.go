package menu

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/menu"
)

func ProvideHandler(
	h handler.HttpHandler,
	srv *menu.Service,
) *Handler {
	return &Handler{
		handler.DefaultCrudHandler[*models.Menu]{
			H:   h,
			Srv: srv,
		},
	}
}

type Handler struct {
	handler.DefaultCrudHandler[*models.Menu]
}
