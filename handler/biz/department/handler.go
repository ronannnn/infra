package department

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/department"
)

func ProvideHandler(
	h handler.HttpHandler,
	srv *department.Service,
) *Handler {
	return &Handler{
		handler.DefaultCrudHandler[*models.Department]{
			H:   h,
			Srv: srv,
		},
	}
}

type Handler struct {
	handler.DefaultCrudHandler[*models.Department]
}
