package models

type Menu struct {
	Base
	Type           *MenuType `json:"type" gorm:"size:1"`
	ParentId       *uint     `json:"parentId"`
	Name           *string   `json:"name"`
	I18nKey        *string   `json:"i18nKey"`
	StaticRouteKey *string   `json:"staticRouteKey" gorm:"type:text"` // 静态路由的Key
	Permission     *string   `json:"permission" gorm:"type:text"`
	Order          *string   `json:"order"`
	Apis           *[]Api    `json:"apis" gorm:"many2many:menu_apis"`
}

type MenuType int

const (
	MenuTypeCatalog MenuType = iota
	MenuTypeMenu
	MenuTypeBtn
)

func (Menu) TableName() string {
	return "departments"
}

func (m Menu) FieldColMapper() map[string]string {
	return CamelToSnakeFromStruct(m)
}
