package login

import "github.com/ronannnn/infra/model"

var (
	LOGIN_TYPE_USERNAME_PASSWORD = "username_password"
)

type UsernamePasswordLoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
