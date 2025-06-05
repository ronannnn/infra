package jobtitle

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/jobtitle"
)

func ProvideHandler(
	h handler.HttpHandler,
	srv *jobtitle.Service,
) *Handler {
	return &Handler{
		handler.DefaultCrudHandler[*models.JobTitle]{
			H:   h,
			Srv: srv,
		},
	}
}

type Handler struct {
	handler.DefaultCrudHandler[*models.JobTitle]
}
