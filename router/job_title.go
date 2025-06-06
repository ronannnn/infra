package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewJobTitleRouter(
	handler *biz.JobTitleHandler,
) *JobTitleRouter {
	return &JobTitleRouter{
		DefaultCrudRouter[*model.JobTitle]{
			BasePath: "/job-titles",
			Handler:  handler,
		},
	}
}

type JobTitleRouter struct {
	DefaultCrudRouter[*model.JobTitle]
}
