package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewJobGradeRouter(
	handler *biz.JobGradeHandler,
) *JobGradeRouter {
	return &JobGradeRouter{
		DefaultCrudRouter[model.JobGrade]{
			BasePath: "/job-grades",
			Handler:  handler,
		},
	}
}

type JobGradeRouter struct {
	DefaultCrudRouter[model.JobGrade]
}
