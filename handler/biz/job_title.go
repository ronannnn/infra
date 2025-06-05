package biz

import (
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
)

func NewJobTitleHandler(
	h handler.HttpHandler,
	srv *service.JobTitleService,
) *JobTitleHandler {
	return &JobTitleHandler{
		handler.DefaultCrudHandler[*model.JobTitle]{
			H:   h,
			Srv: srv,
		},
	}
}

type JobTitleHandler struct {
	handler.DefaultCrudHandler[*model.JobTitle]
}
