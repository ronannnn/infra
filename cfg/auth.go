package cfg

type Auth struct {
	AccessTokenHourDuration  int    `mapstructure:"access-token-hour-duration"`  // Access Token有效时长
	RefreshTokenHourDuration int    `mapstructure:"refresh-token-hour-duration"` // Refresh Token有效时长
	AccessTokenSecret        string `mapstructure:"access-token-secret"`         // Access Token加密字符串
	RefreshTokenSecret       string `mapstructure:"refresh-token-secret"`        // Refresh Token加密字符串
	Enabled                  bool   `mapstructure:"enabled"`                     // 是否启用jwt
}
