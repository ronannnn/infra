package company

import "github.com/ronannnn/infra/models"

type Company struct {
	models.Base
	Name string `json:"name"`
}

type CompanyQuery struct {
	WhereQuery CompanyWhereQuery `json:"whereQuery"`
	OrderQuery CompanyOrderQuery `json:"orderQuery"`
}

type CompanyWhereQuery struct {
	Name string `json:"name" search:"type:like;column:name"`
}

type CompanyOrderQuery struct {
	CreatedAt string `json:"createdAt" search:"type:order;column:created_at"`
}
