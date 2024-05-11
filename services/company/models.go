package company

import "github.com/ronannnn/infra/models"

type Company struct {
	models.Base
	Name string `json:"name"`
}

func (Company) TableName() string {
	return "companies"
}

func (c Company) FieldColMapper() map[string]string {
	return models.CamelToSnakeFromStruct(c)
}
