package router

import (
	"github.com/ronannnn/infra/handler/biz"
	"github.com/ronannnn/infra/model"
)

func NewLoginRecordRouter(
	handler *biz.LoginRecordHandler,
) *LoginRecordRouter {
	return &LoginRecordRouter{
		DefaultCrudRouter[*model.LoginRecord]{
			Handler: handler,
		},
	}
}

type LoginRecordRouter struct {
	DefaultCrudRouter[*model.LoginRecord]
}

func (c *LoginRecordRouter) GetBasePath() string {
	return "/login-records"
}
