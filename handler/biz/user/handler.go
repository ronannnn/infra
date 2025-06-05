package user

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/user"
)

func ProvideHandler(
	h handler.HttpHandler,
	srv *user.Service,
) *Handler {
	return &Handler{
		handler.DefaultCrudHandler[*models.User]{
			H:   h,
			Srv: srv,
		},
	}
}

type Handler struct {
	handler.DefaultCrudHandler[*models.User]
}
