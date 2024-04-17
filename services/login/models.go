package login

type Result struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
}
