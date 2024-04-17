package api

import "github.com/ronannnn/infra/models/request/query"

type ApiQuery struct {
	Pagination  query.Pagination `json:"pagination" search:"-" excel:"-"`
	WhereQuery  ApiWhereQuery    `json:"whereQuery" excel:"-"`
	OrderQuery  ApiOrderQuery    `json:"orderQuery" excel:"-"`
	SelectQUery ApiSelectQuery   `json:"selectQUery"`
}

type ApiWhereQuery struct {
	Name        string   `json:"name" search:"type:like;column:name"`
	Method      []string `json:"method" search:"type:in;column:method"`
	Path        string   `json:"path" search:"type:like;column:path"`
	Description string   `json:"description" search:"type:like;column:description"`
}

type ApiOrderQuery struct {
	CreatedAt   string `json:"createdAt" search:"type:order;column:created_at"`
	Name        string `json:"name" search:"type:order;column:name"`
	Method      string `json:"method" search:"type:order;column:method"`
	Path        string `json:"path" search:"type:order;column:path"`
	Description string `json:"description" search:"type:order;column:description"`
}

type ApiSelectQuery struct {
	Name        bool `json:"name" search:"type:select;column:name" excel:"column:Name;width:20"`
	Method      bool `json:"method" search:"type:select;column:method" excel:"column:Method;width:10"`
	Path        bool `json:"path" search:"type:select;column:path" excel:"column:Path;width:20"`
	Description bool `json:"description" search:"type:select;column:description" excel:"column:Description;width:20"`
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
