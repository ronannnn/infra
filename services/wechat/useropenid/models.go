package useropenid

type UserOpenId struct {
	UserId               uint   `json:"userId" gorm:"user_app_unique_index"`
	OfficialAccountAppId string `json:"officialAccountAppId" gorm:"user_app_unique_index"`
	OpenId               string `json:"openId"`
}
