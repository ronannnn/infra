package loginrecord

import (
	"time"

	"github.com/ronannnn/infra/models"
)

type LoginRecord struct {
	models.Base
	UserId    uint      `json:"userId"`
	DeviceId  string    `json:"deviceId"` // 前端生成的UUID
	LoginTime time.Time `json:"loginTime"`
	Ip        string    `json:"ip"`
	UserAgent string    `json:"userAgent"`
}
