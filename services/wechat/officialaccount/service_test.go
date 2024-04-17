package officialaccount_test

import (
	"testing"

	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/services/wechat/officialaccount"
	"github.com/stretchr/testify/require"
)

func TestWeChatSendTemplateMessage(t *testing.T) {
	var err error
	testCfg := struct {
		WechatOfficialAccount cfg.WechatOfficialAccount `mapstructure:"wechat-official-account"`
	}{}
	err = cfg.ReadFromFile("../../../configs/config.wechattest.toml", &testCfg)
	require.NoError(t, err)

	srv := officialaccount.ProvideService(&testCfg.WechatOfficialAccount)
	err = srv.RefreshAccessToken()
	require.NoError(t, err)
	result, err := srv.SendTemplateMessage(officialaccount.SendTemplateMessageCmd{
		ToUser:     "o0WeS6Sphhj6ZmaDcU-R_0KVKNOo",
		TemplateId: "-tlxQybpqXy69okYsapTjNp4GXI3mdHjGsNE7HlcsJc",
		Url:        "https://photo.gx-logistics.com/tasks/common",
		Data: map[string]map[string]string{
			"inOrderId": {
				"value": "BSX24040009",
			},
			"warehouseCode": {
				"value": "1104-A30",
			},
			"commodity": {
				"value": "混合橡胶 3L",
			},
			"blNumber": {
				"value": "A59D028343",
			},
			"companyName": {
				"value": "广州利泽化工有限公司",
			},
			"creator": {
				"value": "张三",
			},
			"createdAt": {
				"value": "2021-08-24 10:00:00",
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, 0, result.ErrCode)
}

// http://open.weixin.qq.com/connect/oauth2/authorize?appid=wxd2dd67917191fb08&response_type=code&scope=snsapi_base&state=STATE&redirect_uri=http://crmweixin.gx-logistics.com/webMobile/Gate/WebMobileLogin
// http://open.weixin.qq.com/connect/oauth2/authorize?appid=wxd2dd67917191fb08&response_type=code&scope=snsapi_base&state=STATE&redirect_uri=http://crmweixin.gx-logistics.com/webMobile/Gate/WebMobileLogin
// http://open.weixin.qq.com/connect/oauth2/authorize?appid=wxd2dd67917191fb08&response_type=code&scope=snsapi_base&state=STATE&redirect_uri=https://photo.gx-logistics.com/tasks/in

// https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxd2dd67917191fb08&response_type=code&scope=snsapi_base&state=STATE&redirect_uri=https://photo.gx-logistics.com/tasks/in#wechat_redirect