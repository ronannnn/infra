package rowrecord

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/rowrecord"
)

func ProvideHandler(
	h handler.HttpHandler,
	srv *rowrecord.Service,
) *Handler {
	return &Handler{
		handler.DefaultCrudHandler[*models.RowRecord]{
			H:   h,
			Srv: srv,
		},
	}
}

type Handler struct {
	handler.DefaultCrudHandler[*models.RowRecord]
}
