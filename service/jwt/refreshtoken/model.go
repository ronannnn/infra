package refreshtoken

import (
	"net/http"

	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/utils/useragent"
)

// RefreshToken 记录用户的Refresh Token，控制登录设备数量
type RefreshToken struct {
	model.Base
	UserId          *uint                 `json:"userId" gorm:"uniqueIndex:user_id_device_refresh_token_idx"`
	LoginDeviceType *useragent.DeviceType `json:"loginDeviceType" gorm:"uniqueIndex:user_id_device_refresh_token_idx"`
	RefreshToken    *string               `json:"refreshToken"`
	DeviceId        *string               `json:"deviceId"` // 前端生成的UUID
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (refreshToken RefreshToken) WithOprFromReq(r *http.Request) model.Crudable {
	refreshToken.OprBy = model.GetOprFromReq(r)
	return refreshToken
}

func (refreshToken RefreshToken) WithUpdaterFromReq(r *http.Request) model.Crudable {
	refreshToken.OprBy = model.GetUpdaterFromReq(r)
	return refreshToken
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
