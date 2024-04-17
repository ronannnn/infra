package api

import "github.com/ronannnn/infra/models"

type Api struct {
	models.Base
	Name        *string `json:"name"`
	Method      *string `json:"method"`
	Path        *string `json:"path"`
	Description *string `json:"description"`
}
