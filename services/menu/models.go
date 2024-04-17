package menu

import "github.com/ronannnn/infra/models"

type Menu struct {
	models.Base
	Type           *MenuType     `json:"type" gorm:"size:1"`
	ParentId       *uint         `json:"parentId"`
	Name           *string       `json:"name"`
	I18nKey        *string       `json:"i18nKey"`
	StaticRouteKey *string       `json:"staticRouteKey" gorm:"type:text"` // 静态路由的Key
	Permission     *string       `json:"permission" gorm:"type:text"`
	Order          *int          `json:"order"`
	Apis           *[]models.Api `json:"apis" gorm:"many2many:menu_apis"`
}

type MenuType int

const (
	MenuTypeCatalog MenuType = iota
	MenuTypeMenu
	MenuTypeBtn
)

type MenuQuery struct {
	WhereQuery MenuWhereQuery `json:"whereQuery"`
	OrderQuery MenuOrderQuery `json:"orderQuery"`
}

type MenuWhereQuery struct {
	Type           []MenuType `json:"type" search:"type:in;column:type"`
	Name           string     `json:"name" search:"type:like;column:name"`
	StaticRouteKey string     `json:"staticRouteKey" search:"type:like;column:static_route_key"`
	Permission     string     `json:"permission" search:"type:like;column:permission"`
}

type MenuOrderQuery struct {
	CreatedAt string `json:"createdAt" search:"type:order;column:created_at"`
}
