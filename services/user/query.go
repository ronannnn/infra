package user

import "github.com/ronannnn/infra/models/request/query"

type UserQuery struct {
	Pagination query.Pagination `json:"pagination"`
	WhereQuery UserWhereQuery   `json:"whereQuery" query:"category:where"`
	OrderQuery []UserOrderQuery `json:"orderQuery" query:"category:order"`
}

type UserWhereQuery struct {
	Nickname string `json:"nickname" query:"type:like;column:nickname"`
	Username string `json:"username" query:"type:like;column:username"`
	Email    string `json:"email" query:"type:like;column:email"`
	TelNo    string `json:"telNo" query:"type:like;column:telNo"`
}

type UserOrderQuery struct {
	CreatedAt string `json:"createdAt" query:"column:created_at"`
	Nickname  string `json:"nickname" query:"column:nickname"`
	Username  string `json:"username" query:"column:username"`
	Email     string `json:"email" query:"column:email"`
	TelNo     string `json:"telNo" query:"column:telNo"`
}
