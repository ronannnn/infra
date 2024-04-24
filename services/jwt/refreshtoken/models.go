package refreshtoken

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/utils/useragent"
)

// RefreshToken 记录用户的Refresh Token，控制登录设备数量
type RefreshToken struct {
	models.Base
	UserId          *uint                 `json:"userId" gorm:"uniqueIndex:user_id_device_refresh_token_idx"`
	LoginDeviceType *useragent.DeviceType `json:"loginDeviceType" gorm:"uniqueIndex:user_id_device_refresh_token_idx"`
	RefreshToken    *string               `json:"refreshToken"`
	DeviceId        *string               `json:"deviceId"` // 前端生成的UUID
}

type BaseClaims struct {
	UserId   uint   `json:"userId"` // The first letter must be upper case, or parseToken() cannot get user id
	Username string `json:"username"`
}

func (c *BaseClaims) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"username": c.Username,
		"userId":   c.UserId,
	}
}
