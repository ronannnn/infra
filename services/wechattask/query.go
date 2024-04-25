package wechattask

import "github.com/ronannnn/infra/models/request/query"

type WechatTaskQuery struct {
	Pagination query.Pagination     `json:"pagination" search:"-"`
	WhereQuery WechatTaskWhereQuery `json:"whereQuery"`
	OrderQuery WechatTaskOrderQuery `json:"orderQuery"`
}

type WechatTaskWhereQuery struct {
	Name     string `json:"name" search:"type:like;column:name"`
	Uuid     string `json:"uuid" search:"type:like;column:uuid"`
	Disabled *int8  `json:"disabled" search:"type:equal;column:disabled"`
}

type WechatTaskOrderQuery struct {
	CreatedAt string `json:"createdAt" search:"type:order;column:created_at"`
}
