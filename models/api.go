package models

type Api struct {
	Base
	Name        *string `json:"name"`
	Method      *string `json:"method"`
	Path        *string `json:"path"`
	Description *string `json:"description"`
}
