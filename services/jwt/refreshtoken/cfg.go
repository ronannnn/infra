package refreshtoken

type Cfg struct {
	RefreshTokenMinuteDuration int    `mapstructure:"refresh-token-minute-duration"` // Refresh Token有效时长
	RefreshTokenSecret         string `mapstructure:"refresh-token-secret"`          // Refresh Token加密字符串
}
