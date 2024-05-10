package wechattask

import "github.com/ronannnn/infra/models/request/query"

type WechatTaskQuery struct {
	Pagination query.Pagination       `json:"pagination"`
	WhereQuery WechatTaskWhereQuery   `json:"whereQuery" query:"category:where"`
	OrderQuery []WechatTaskOrderQuery `json:"orderQuery" query:"category:order"`
}

type WechatTaskWhereQuery struct {
	Name     string `json:"name" query:"type:like;column:name"`
	Uuid     string `json:"uuid" query:"type:like;column:uuid"`
	Disabled *int8  `json:"disabled" query:"type:equal;column:disabled"`
}

type WechatTaskOrderQuery struct {
	CreatedAt string `json:"createdAt" query:"column:created_at"`
}
