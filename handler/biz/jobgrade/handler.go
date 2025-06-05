package jobgrade

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/jobgrade"
)

func ProvideHandler(
	h handler.HttpHandler,
	srv *jobgrade.Service,
) *Handler {
	return &Handler{
		handler.DefaultCrudHandler[*models.JobGrade]{
			H:   h,
			Srv: srv,
		},
	}
}

type Handler struct {
	handler.DefaultCrudHandler[*models.JobGrade]
}
