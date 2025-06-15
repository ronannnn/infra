package model

import "net/http"

type Department struct {
	Base
	Name      *string     `json:"name"`
	CompanyId *uint       `json:"companyId"`
	Company   *Company    `json:"company"`
	LeaderId  *uint       `json:"leaderId"`
	Leader    *User       `json:"leader"`
	ParentId  *uint       `json:"parentId"`
	Parent    *Department `json:"parent"`
}

func (Department) TableName() string {
	return "departments"
}

func (model Department) WithOprFromReq(r *http.Request) Crudable {
	model.OprBy = GetOprFromReq(r)
	return model
}

func (model Department) WithUpdaterFromReq(r *http.Request) Crudable {
	model.OprBy = GetUpdaterFromReq(r)
	return model
}
