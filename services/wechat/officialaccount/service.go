package officialaccount

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ronannnn/infra/cfg"
)

type Service interface {
	GetAccessToken() (GetAccessTokenResult, error)
	RefreshAccessToken() error

	GetJsApiTicket() (GetJsApiTicketResult, error)
	RefreshJsApiTicket() error
	GetSignedJsSdkConfig(GetJsSdkConfigCmd) (GetJsSdkConfigResult, error)

	GetOpenId(GetOpenIdCmd) (GetOpenIdResult, error)

	SendTemplateMessage(SendTemplateMessageCmd) (SendTemplateMessageResult, error)

	SendSubscriptionTemplateMessage(SendSubscriptionTemplateMessageCmd) (SendSubscriptionTemplateMessageResult, error)
}

type ServiceImpl struct {
	cfg         *cfg.WechatOfficialAccount
	AccessToken string
	JsApiTicket string
}

func ProvideService(cfg *cfg.WechatOfficialAccount) Service {
	return &ServiceImpl{cfg: cfg}
}

// GetAccessToken 获取 access_token
// 文档地址：https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
// 接口：   GET https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
func (s *ServiceImpl) GetAccessToken() (result GetAccessTokenResult, err error) {
	params := url.Values{
		"appid":      {s.cfg.AppId},
		"secret":     {s.cfg.AppSecret},
		"grant_type": {"client_credential"},
	}
	var resp *http.Response
	if resp, err = http.Get("https://api.weixin.qq.com/cgi-bin/token?" + params.Encode()); err != nil {
		return
	}
	defer resp.Body.Close()
	// result
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetAccessToken failed: %s", result.ErrMsg)
	}
	return
}

func (s *ServiceImpl) RefreshAccessToken() (err error) {
	var accessTokenResult GetAccessTokenResult
	if accessTokenResult, err = s.GetAccessToken(); err != nil {
		return
	}
	s.AccessToken = accessTokenResult.AccessToken
	return
}

// GetJsApiTicket 获取 jsapi_ticket
// 文档地址：https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html
// 接口：   GET https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=ACCESS_TOKEN&type=jsapi
func (s *ServiceImpl) GetJsApiTicket() (result GetJsApiTicketResult, err error) {
	if s.AccessToken == "" {
		if err = s.RefreshAccessToken(); err != nil {
			return
		}
	}
	params := url.Values{
		"access_token": {s.AccessToken},
		"type":         {"jsapi"},
	}
	var resp *http.Response
	if resp, err = http.Get("https://api.weixin.qq.com/cgi-bin/ticket/getticket?" + params.Encode()); err != nil {
		return
	}
	defer resp.Body.Close()
	// result
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetJsApiTicket failed: %s", result.ErrMsg)
	}
	return
}

func (s *ServiceImpl) RefreshJsApiTicket() (err error) {
	var jsApiTicketResult GetJsApiTicketResult
	if jsApiTicketResult, err = s.GetJsApiTicket(); err != nil {
		return
	}
	s.JsApiTicket = jsApiTicketResult.Ticket
	return
}

// GetSignedJsSdkConfig 获取签名后的 js-sdk 配置
func (s *ServiceImpl) GetSignedJsSdkConfig(cmd GetJsSdkConfigCmd) (result GetJsSdkConfigResult, err error) {
	if s.JsApiTicket == "" {
		if err = s.RefreshJsApiTicket(); err != nil {
			return
		}
	}
	result.JsApiTicket = s.JsApiTicket
	result.NonceStr = cmd.NonceStr
	result.Timestamp = time.Now().Unix()
	result.Url = cmd.Url
	result.Signature = signJsSdkConfig(s.JsApiTicket, result.NonceStr, result.Timestamp, result.Url)
	return
}

// GetOpenId 获取 openid
// 文档地址：https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html
// 接口：   GET https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code
func (s *ServiceImpl) GetOpenId(cmd GetOpenIdCmd) (result GetOpenIdResult, err error) {
	params := url.Values{
		"appid":      {s.cfg.AppId},
		"secret":     {s.cfg.AppSecret},
		"code":       {cmd.Code},
		"grant_type": {"authorization_code"},
	}
	var resp *http.Response
	if resp, err = http.Get("https://api.weixin.qq.com/sns/oauth2/access_token?" + params.Encode()); err != nil {
		return
	}
	defer resp.Body.Close()
	// result
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("GetOpenId failed: %s", result.ErrMsg)
	}
	return
}

func (s *ServiceImpl) SendTemplateMessage(cmd SendTemplateMessageCmd) (result SendTemplateMessageResult, err error) {
	if s.AccessToken == "" {
		if err = s.RefreshAccessToken(); err != nil {
			return
		}
	}
	if result, err = s.sendTemplateMessageFn(cmd); err != nil {
		if err = s.RefreshAccessToken(); err != nil {
			return
		}
		result, err = s.sendTemplateMessageFn(cmd)
	}
	return
}

// SendTemplateMessage 发送模板消息
// 文档地址：https://developers.weixin.qq.com/doc/offiaccount/Message_Management/Template_Message_Interface.html
// 接口：   POST https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=ACCESS_TOKEN
func (s *ServiceImpl) sendTemplateMessageFn(cmd SendTemplateMessageCmd) (result SendTemplateMessageResult, err error) {
	// req
	params := url.Values{
		"access_token": {s.AccessToken},
	}
	var jsonData []byte
	if jsonData, err = json.Marshal(cmd); err != nil {
		return
	}
	var resp *http.Response
	if resp, err = http.Post("https://api.weixin.qq.com/cgi-bin/message/template/send?"+params.Encode(), "application/json", bytes.NewBuffer(jsonData)); err != nil {
		return
	}
	defer resp.Body.Close()
	// result
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("SendTemplateMessage failed: %s", result.ErrMsg)
	}
	return
}

func (s *ServiceImpl) SendSubscriptionTemplateMessage(cmd SendSubscriptionTemplateMessageCmd) (result SendSubscriptionTemplateMessageResult, err error) {
	if s.AccessToken == "" {
		if err = s.RefreshAccessToken(); err != nil {
			return
		}
	}
	if result, err = s.sendSubscriptionTemplateMessageFn(cmd); err != nil {
		if err = s.RefreshAccessToken(); err != nil {
			return
		}
		result, err = s.sendSubscriptionTemplateMessageFn(cmd)
	}
	return
}

// SendSubscriptionTemplateMessage 发送模板消息
// 文档地址：https://developers.weixin.qq.com/doc/offiaccount/Subscription_Messages/api.html
// 接口：   POST https://api.weixin.qq.com/cgi-bin/message/subscribe/bizsend?access_token=ACCESS_TOKEN
func (s *ServiceImpl) sendSubscriptionTemplateMessageFn(cmd SendSubscriptionTemplateMessageCmd) (result SendSubscriptionTemplateMessageResult, err error) {
	// req
	params := url.Values{
		"access_token": {s.AccessToken},
	}
	var jsonData []byte
	if jsonData, err = json.Marshal(cmd); err != nil {
		return
	}
	var resp *http.Response
	if resp, err = http.Post("https://api.weixin.qq.com/cgi-bin/message/subscribe/bizsend?"+params.Encode(), "application/json", bytes.NewBuffer(jsonData)); err != nil {
		return
	}
	defer resp.Body.Close()
	// result
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}
	if result.ErrCode != 0 {
		err = fmt.Errorf("SendTemplateMessage failed: %s", result.ErrMsg)
	}
	return
}
