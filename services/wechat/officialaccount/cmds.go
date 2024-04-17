package officialaccount

type GetAccessTokenResult struct {
	CommonResult
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type GetOpenIdCmd struct {
	Code string `json:"code"`
}

type GetOpenIdResult struct {
	CommonResult
	OpenId         string `json:"openid"`
	AccessToken    string `json:"access_token"`
	ExpiresIn      int    `json:"expires_in"`
	RefreshToken   string `json:"refresh_token"`
	Scope          string `json:"scope"`
	UnionId        string `json:"unionid"`
	IsSnapshotUser int    `json:"is_snapshotuser"`
}

type SendTemplateMessageCmd struct {
	ToUser       string                       `json:"touser"`
	TemplateId   string                       `json:"template_id"`
	Url          string                       `json:"url"`
	Data         map[string]map[string]string `json:"data"`
	ClientMsgsId string                       `json:"client_msg_id"`
}

type SendTemplateMessageResult struct {
	CommonResult
	MsgId int64 `json:"msgid"`
}

type CommonResult struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
