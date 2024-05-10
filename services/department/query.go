package department

import "github.com/ronannnn/infra/models/request/query"

type DepartmentQuery struct {
	Pagination query.Pagination       `json:"pagination"`
	WhereQuery []DepartmentWhereQuery `json:"whereQuery" query:"category:where"`
	OrderQuery []DepartmentOrderQuery `json:"orderQuery" query:"category:order"`
}

type DepartmentWhereQuery struct {
	Name string `json:"name" query:"type:like;column:name"`
}

type DepartmentOrderQuery struct {
	CreatedAt string `json:"createdAt" query:"column:created_at"`
}
