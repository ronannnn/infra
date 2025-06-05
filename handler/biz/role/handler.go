package role

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/role"
)

func ProvideHandler(
	h handler.HttpHandler,
	srv *role.Service,
) *Handler {
	return &Handler{
		handler.DefaultCrudHandler[*models.Role]{
			H:   h,
			Srv: srv,
		},
	}
}

type Handler struct {
	handler.DefaultCrudHandler[*models.Role]
}
