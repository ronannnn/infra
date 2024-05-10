package company

import "github.com/ronannnn/infra/models/request/query"

type CompanyQuery struct {
	Pagination query.Pagination    `json:"pagination"`
	WhereQuery []CompanyWhereQuery `json:"whereQuery" query:"category:where"`
	OrderQuery []CompanyOrderQuery `json:"orderQuery" query:"category:order"`
}

type CompanyWhereQuery struct {
	Name string `json:"name" query:"type:like;column:name"`
}

type CompanyOrderQuery struct {
	CreatedAt string `json:"createdAt" query:"column:created_at"`
}
