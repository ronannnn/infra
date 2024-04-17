package cfg

type WechatOfficialAccount struct {
	AppId     string `mapstructure:"app-id"`
	AppSecret string `mapstructure:"app-secret"`
}
