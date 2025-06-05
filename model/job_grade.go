package model

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
