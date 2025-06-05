package biz

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewJobGradeHandler(
	h handler.HttpHandler,
	srv *service.JobGradeService,
) *JobGradeHandler {
	return &JobGradeHandler{
		handler.DefaultCrudHandler[*model.JobGrade]{
			H:   h,
			Srv: srv,
		},
	}
}

type JobGradeHandler struct {
	handler.DefaultCrudHandler[*model.JobGrade]
}
