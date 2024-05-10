package api

import "github.com/ronannnn/infra/models/request/query"

type ApiQuery struct {
	Pagination  query.Pagination `json:"pagination" excel:"-"`
	WhereQuery  []ApiWhereQuery  `json:"whereQuery" query:"category:where" excel:"-"`
	OrderQuery  []ApiOrderQuery  `json:"orderQuery" query:"category:order" excel:"-"`
	SelectQUery []ApiSelectQuery `json:"selectQUery" query:"category:select"`
}

type ApiWhereQuery struct {
	Name        string   `json:"name" query:"type:like;column:name"`
	Method      []string `json:"method" query:"type:in;column:method"`
	Path        string   `json:"path" query:"type:like;column:path"`
	Description string   `json:"description" query:"type:like;column:description"`
}

type ApiOrderQuery struct {
	CreatedAt   string `json:"createdAt" query:"column:created_at"`
	Name        string `json:"name" query:"column:name"`
	Method      string `json:"method" query:"column:method"`
	Path        string `json:"path" query:"column:path"`
	Description string `json:"description" query:"column:description"`
}

type ApiSelectQuery struct {
	Name        bool `json:"name" query:"column:name" excel:"column:Name;width:20"`
	Method      bool `json:"method" query:"column:method" excel:"column:Method;width:10"`
	Path        bool `json:"path" query:"column:path" excel:"column:Path;width:20"`
	Description bool `json:"description" query:"column:description" excel:"column:Description;width:20"`
}

func (a *ApiQuery) TableName() string {
	return "apis"
}

func (a *ApiQuery) SheetName() string {
	return "apis"
}

func (a *ApiQuery) Model() any {
	return *a
}
