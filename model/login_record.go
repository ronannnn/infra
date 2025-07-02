package model

import (
	"net/http"
	"time"

	"github.com/ronannnn/infra/utils/useragent"
)

type LoginStatus int

const (
	LoginStatusSuccess LoginStatus = iota + 1
	LoginStatusFailed
	LoginStatusDupLogin
	LoginStatusErrUsernameOrPassword
	LoginStatusErrUserNotExists
)

type LoginRecord struct {
	Base
	UserId          *uint                 `json:"userId"`
	DeviceId        *string               `json:"deviceId"` // 前端生成的UUID
	LoginDeviceType *useragent.DeviceType `json:"loginDeviceType"`
	LoginTime       *time.Time            `json:"loginTime"`
	Ip              *string               `json:"ip"`
	UserAgent       *string               `json:"userAgent"`
	Status          *LoginStatus          `json:"status"`
	LoginType       *string               `json:"loginType"` // 对应User.LoginType
}

func (LoginRecord) TableName() string {
	return "login_records"
}

func (model LoginRecord) WithOprFromReq(r *http.Request) Crudable {
	model.OprBy = GetOprFromReq(r)
	return model
}

func (model LoginRecord) WithUpdaterFromReq(r *http.Request) Crudable {
	model.OprBy = GetUpdaterFromReq(r)
	return model
}
