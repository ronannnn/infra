package role

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
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

type RoleQuery struct {
	Pagination query.Pagination `json:"pagination" search:"-"`
	WhereQuery RoleWhereQuery   `json:"whereQuery"`
	OrderQuery RoleOrderQuery   `json:"orderQuery"`
}

type RoleWhereQuery struct {
	Name       string `json:"name" search:"type:like;column:name"`
	Permission string `json:"permission" search:"type:like;column:description"`
	Disabled   []int  `json:"disabled" search:"type:in;column:disabled"`
	Remark     string `json:"remark" search:"type:like;column:remark"`
}

type RoleOrderQuery struct {
	CreatedAt  string `json:"createdAt" search:"type:order;column:created_at"`
	Name       string `json:"name" search:"type:order;column:name"`
	Permission string `json:"permission" search:"type:order;column:permission"`
	Disabled   bool   `json:"disabled" search:"type:order;column:disabled"`
}
