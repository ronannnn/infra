package model

type JobTitle struct {
	Base
	Name        *string `json:"name"`        // 职务名称
	Description *string `json:"description"` // 职务描述
	Disabled    *bool   `json:"disabled"`    // 是否禁用
	Remark      *string `json:"remark"`      // 备注
}

func (JobTitle) TableName() string {
	return "job_titles"
}
