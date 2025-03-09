package models

type Role struct {
	Base
	Name       *string `json:"name"`       // 角色名称
	Permission *string `json:"permission"` // 角色权限标识符
	Disabled   *bool   `json:"disabled"`   // 角色是否禁用
	Remark     *string `json:"remark"`     // 备注
	Menus      []*Menu `json:"menus" gorm:"many2many:role_menus;"`
}

func (Role) TableName() string {
	return "roles"
}

func (r Role) FieldColMapper() map[string]string {
	return CamelToSnakeFromStruct(r)
}
