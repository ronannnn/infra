package models

type Api struct {
	Base
	Name        *string `json:"name"`
	Method      *string `json:"method"`
	Path        *string `json:"path"`
	Description *string `json:"description"`
}

func (Api) TableName() string {
	return "apis"
}

func (a Api) FieldColMapper() map[string]string {
	return CamelToSnakeFromStruct(a)
}
