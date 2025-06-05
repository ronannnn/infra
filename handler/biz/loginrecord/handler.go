package loginrecord

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/loginrecord"
)

func ProvideHandler(
	h handler.HttpHandler,
	srv *loginrecord.Service,
) *Handler {
	return &Handler{
		handler.DefaultCrudHandler[*models.LoginRecord]{
			H:   h,
			Srv: srv,
		},
	}
}

type Handler struct {
	handler.DefaultCrudHandler[*models.LoginRecord]
}
