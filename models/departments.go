package models

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
