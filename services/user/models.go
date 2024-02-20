package user

// request commands and response results
// login with username and password
type LoginUsernameCommand struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshTokensCommand struct {
	RefreshToken string `json:"refreshToken"`
}

type AuthResult struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
}

type ChangeUserLoginPwdCommand struct {
	UserId uint   `json:"userId"`
	OldPwd string `json:"oldPwd"`
	NewPwd string `json:"newPwd"`
}
