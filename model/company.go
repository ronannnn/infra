package model

type Company struct {
	Base
	Name string `json:"name"`
}

func (Company) TableName() string {
	return "companies"
}
