package company

import "github.com/ronannnn/infra/models"

type Company struct {
	models.Base
	Name string `json:"name"`
}
