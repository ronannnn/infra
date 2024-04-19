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
	err = cfg.ReadFromFile("../../../configs/config.wechattest1.toml", &testCfg)
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

func TestWeChatSendSubscriptionTemplateMessage(t *testing.T) {
	var err error
	testCfg := struct {
		WechatOfficialAccount cfg.WechatOfficialAccount `mapstructure:"wechat-official-account"`
	}{}
	err = cfg.ReadFromFile("../../../configs/config.wechattest2.toml", &testCfg)
	require.NoError(t, err)

	srv := officialaccount.ProvideService(&testCfg.WechatOfficialAccount)
	err = srv.RefreshAccessToken()
	require.NoError(t, err)
	result, err := srv.SendTemplateMessage(officialaccount.SendTemplateMessageCmd{
		ToUser:     "o0WeS6Sphhj6ZmaDcU-R_0KVKNOo",
		TemplateId: "tOAHBHlWhVckcpc1jzm03qcHxJy2rxlfMCmbzxBuVfw",
		Url:        "https://photo.gx-logistics.com/tasks/common",
		Data: map[string]map[string]string{
			"character_string1": {
				"value": "BSX24040009",
			},
			"character_string2": {
				"value": "1104-A30",
			},
			"thing3": {
				"value": "混合橡胶 3L 30吨 900件",
			},
			"character_string4": {
				"value": "A59D028343",
			},
			"thing5": {
				"value": "拍一下粒子照片",
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, 0, result.ErrCode)
}
