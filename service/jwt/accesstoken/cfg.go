package accesstoken

type Cfg struct {
	AccessTokenMinuteDuration int    `mapstructure:"access-token-minute-duration"` // Access Token有效时长
	AccessTokenSecret         string `mapstructure:"access-token-secret"`          // Access Token加密字符串
	Enabled                   bool   `mapstructure:"enabled"`                      // 是否启用jwt
}
