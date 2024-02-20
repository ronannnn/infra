package query

import "github.com/ronannnn/infra/models"

type MenuQuery struct {
	WhereQuery MenuWhereQuery `json:"whereQuery"`
	OrderQuery MenuOrderQuery `json:"orderQuery"`
}

type MenuWhereQuery struct {
	Type           []models.MenuType `json:"type" search:"type:in;column:type"`
	Name           string            `json:"name" search:"type:like;column:name"`
	StaticRouteKey string            `json:"staticRouteKey" search:"type:like;column:static_route_key"`
	Permission     string            `json:"permission" search:"type:like;column:permission"`
}

type MenuOrderQuery struct {
	CreatedAt string `json:"createdAt" search:"type:order;column:created_at"`
}
