package login

import "github.com/ronannnn/infra/model"

// request commands and response results
// login with username and password
type UsernameCmd struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserAgent string `json:"userAgent"`
	DeviceId  string `json:"deviceId"`
}

type Result struct {
	RefreshToken string      `json:"refreshToken"`
	AccessToken  string      `json:"accessToken"`
	DupLogin     bool        `json:"dupLogin"`
	User         *model.User `json:"user"`
}

// refresh refresh token and access token
type RefreshTokensCmd struct {
	RefreshToken string `json:"refreshToken"`
}

type ChangeUserPwdCmd struct {
	UserId uint   `json:"userId"`
	OldPwd string `json:"oldPwd"`
	NewPwd string `json:"newPwd"`
}
