package loginrecord

import (
	"time"

	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/utils/useragent"
)

type Status int

const (
	StatusSuccess Status = iota + 1
	StatusFailed
	StatusDupLogin
	StatusErrUsernameOrPassword
	StatusErrUserNotExists
)

type LoginRecord struct {
	models.Base
	UserId          uint                  `json:"userId"`
	DeviceId        string                `json:"deviceId"` // 前端生成的UUID
	LoginDeviceType *useragent.DeviceType `json:"loginDeviceType"`
	LoginTime       time.Time             `json:"loginTime"`
	Ip              string                `json:"ip"`
	UserAgent       string                `json:"userAgent"`
	Status          Status                `json:"status"`
}
