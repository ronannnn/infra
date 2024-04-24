package login

// request commands and response results
// login with username and password
type UsernameCmd struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserAgent string `json:"userAgent"`
	DeviceId  string `json:"deviceId"`
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
