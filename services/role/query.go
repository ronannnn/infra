package role

import "github.com/ronannnn/infra/models/request/query"

type RoleQuery struct {
	Pagination query.Pagination `json:"pagination"`
	WhereQuery RoleWhereQuery   `json:"whereQuery" query:"category:where"`
	OrderQuery []RoleOrderQuery `json:"orderQuery" query:"category:order"`
}

type RoleWhereQuery struct {
	Name       string `json:"name" query:"type:like;column:name"`
	Permission string `json:"permission" query:"type:like;column:description"`
	Disabled   []int  `json:"disabled" query:"type:in;column:disabled"`
	Remark     string `json:"remark" query:"type:like;column:remark"`
}

type RoleOrderQuery struct {
	CreatedAt  string `json:"createdAt" query:"column:created_at"`
	Name       string `json:"name" query:"column:name"`
	Permission string `json:"permission" query:"column:permission"`
	Disabled   bool   `json:"disabled" query:"column:disabled"`
}
