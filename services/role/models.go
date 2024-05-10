package role

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/menu"
)

type Role struct {
	models.Base
	Name       *string      `json:"name"`       // 角色名称
	Permission *string      `json:"permission"` // 角色权限标识符
	Disabled   *bool        `json:"disabled"`   // 角色是否禁用
	Remark     *string      `json:"remark"`     // 备注
	Menus      *[]menu.Menu `json:"menus" gorm:"many2many:role_menus;"`
}
