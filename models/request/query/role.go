package query

type RoleQuery struct {
	Pagination Pagination     `json:"pagination" search:"-"`
	WhereQuery RoleWhereQuery `json:"whereQuery"`
	OrderQuery RoleOrderQuery `json:"orderQuery"`
}

type RoleWhereQuery struct {
	Name       string `json:"name" search:"type:like;column:name"`
	Permission string `json:"permission" search:"type:like;column:description"`
	Disabled   []int  `json:"disabled" search:"type:in;column:disabled"`
	Remark     string `json:"remark" search:"type:like;column:remark"`
}

type RoleOrderQuery struct {
	CreatedAt  string `json:"createdAt" search:"type:order;column:created_at"`
	Name       string `json:"name" search:"type:order;column:name"`
	Permission string `json:"permission" search:"type:order;column:permission"`
	Disabled   bool   `json:"disabled" search:"type:order;column:disabled"`
}
