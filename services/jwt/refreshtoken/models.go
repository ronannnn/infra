package refreshtoken

import "github.com/ronannnn/infra/models"

type RefreshToken struct {
	models.Base
	UserId       *uint   `gorm:"uniqueIndex:user_id_refresh_token_idx"`
	RefreshToken *string `gorm:"uniqueIndex:user_id_refresh_token_idx"`
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
