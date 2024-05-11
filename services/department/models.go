package department

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/company"
)

type Department struct {
	models.Base
	Name      *string          `json:"name"`
	CompanyId *uint            `json:"companyId"`
	Company   *company.Company `json:"company"`
	LeaderId  *uint            `json:"leaderId"`
	Leader    *models.User     `json:"leader"`
	ParentId  *uint            `json:"parentId"`
	Parent    *Department      `json:"parent"`
}

func (Department) TableName() string {
	return "departments"
}

func (d Department) FieldColMapper() map[string]string {
	return models.CamelToSnakeFromStruct(d)
}
