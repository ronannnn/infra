package model

import "net/http"

type JobGrade struct {
	Base
	Name        *string `json:"name"`        // 职级名称
	Description *string `json:"description"` // 职级描述
	Disabled    *bool   `json:"disabled"`    // 是否禁用
	Remark      *string `json:"remark"`      // 备注
}

func (JobGrade) TableName() string {
	return "job_grades"
}

func (model JobGrade) WithOprFromReq(r *http.Request) Crudable {
	model.OprBy = GetOprFromReq(r)
	return model
}

func (model JobGrade) WithUpdaterFromReq(r *http.Request) Crudable {
	model.OprBy = GetUpdaterFromReq(r)
	return model
}
