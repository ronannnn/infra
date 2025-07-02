package model

import "net/http"

type Company struct {
	Base
	Name *string `json:"name" gorm:"uniqueIndex"`
}

func (Company) TableName() string {
	return "companies"
}

func (model Company) WithOprFromReq(r *http.Request) Crudable {
	model.OprBy = GetOprFromReq(r)
	return model
}

func (model Company) WithUpdaterFromReq(r *http.Request) Crudable {
	model.OprBy = GetUpdaterFromReq(r)
	return model
}
