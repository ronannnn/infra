package user

import "github.com/ronannnn/infra/models/request/query"

type UserQuery struct {
	Pagination query.Pagination `json:"pagination" search:"-"`
	WhereQuery UserWhereQuery   `json:"whereQuery"`
	OrderQuery UserOrderQuery   `json:"orderQuery"`
}

type UserWhereQuery struct {
	Nickname string `json:"nickname" search:"type:like;column:nickname"`
	Username string `json:"username" search:"type:like;column:username"`
	Email    string `json:"email" search:"type:like;column:email"`
	TelNo    string `json:"telNo" search:"type:like;column:telNo"`
}

type UserOrderQuery struct {
	CreatedAt string `json:"createdAt" search:"type:order;column:created_at"`
	Nickname  string `json:"nickname" search:"type:order;column:nickname"`
	Username  string `json:"username" search:"type:order;column:username"`
	Email     string `json:"email" search:"type:order;column:email"`
	TelNo     string `json:"telNo" search:"type:order;column:telNo"`
}
