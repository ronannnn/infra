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
			Handler: handler,
		},
	}
}

type JobTitleRouter struct {
	DefaultCrudRouter[*model.JobTitle]
}

func (c *JobTitleRouter) GetBasePath() string {
	return "/job-titles"
}
