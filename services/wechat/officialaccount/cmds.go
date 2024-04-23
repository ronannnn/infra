package officialaccount

type GetAccessTokenResult struct {
	CommonResult
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type GetJsApiTicketResult struct {
	CommonResult
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
}

type GetJsSdkConfigCmd struct {
	NonceStr string `json:"noncestr"`
	Url      string `json:"url"`
}

type GetJsSdkConfigResult struct {
	JsApiTicket string `json:"jsApiTicket"`
	NonceStr    string `json:"nonceStr"`
	Timestamp   int64  `json:"timestamp"`
	Url         string `json:"url"`
	Signature   string `json:"signature"`
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

type SendSubscriptionTemplateMessageCmd struct {
	ToUser     string                       `json:"touser"`
	TemplateId string                       `json:"template_id"`
	Page       string                       `json:"page"`
	Data       map[string]map[string]string `json:"data"`
}

type SendSubscriptionTemplateMessageResult struct {
	CommonResult
}

type CommonResult struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
