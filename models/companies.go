package models

type Company struct {
	Base
	Name string `json:"name"`
}

func (Company) TableName() string {
	return "companies"
}

func (c Company) FieldColMapper() map[string]string {
	return CamelToSnakeFromStruct(c)
}
