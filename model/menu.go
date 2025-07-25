package model

import "net/http"

type Menu struct {
	Base
	Type           *MenuType `json:"type"`
	ParentId       *uint     `json:"parentId"`
	Name           *string   `json:"name"`
	I18nKey        *string   `json:"i18nKey"`
	StaticRouteKey *string   `json:"staticRouteKey" gorm:"type:text"` // 静态路由的Key
	Permission     *string   `json:"permission" gorm:"type:text"`
	Order          *string   `json:"order"`
	Disabled       *bool     `json:"disabled"` // 前段会显示出这个菜单，但是不可点击
	Apis           []*Api    `json:"apis" gorm:"many2many:menu_apis"`
}

type MenuType string

const (
	MenuTypeCatalog MenuType = "catalog"
	MenuTypeMenu    MenuType = "menu"
	MenuTypeBtn     MenuType = "btn"
)

func (Menu) TableName() string {
	return "menus"
}

func (model Menu) WithOprFromReq(r *http.Request) Crudable {
	model.OprBy = GetOprFromReq(r)
	return model
}

func (model Menu) WithUpdaterFromReq(r *http.Request) Crudable {
	model.OprBy = GetUpdaterFromReq(r)
	return model
}
