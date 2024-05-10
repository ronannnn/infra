package menu

import "github.com/ronannnn/infra/models/request/query"

type MenuQuery struct {
	Pagination query.Pagination `json:"pagination"`
	WhereQuery []MenuWhereQuery `json:"whereQuery" query:"category:where"`
	OrderQuery []MenuOrderQuery `json:"orderQuery" query:"category:order"`
}

type MenuWhereQuery struct {
	Type           []MenuType `json:"type" query:"type:in;column:type"`
	Name           string     `json:"name" query:"type:like;column:name"`
	StaticRouteKey string     `json:"staticRouteKey" query:"type:like;column:static_route_key"`
	Permission     string     `json:"permission" query:"type:like;column:permission"`
}

type MenuOrderQuery struct {
	CreatedAt string `json:"createdAt" query:"column:created_at"`
}
