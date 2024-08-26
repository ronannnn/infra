package redis

type Cfg struct {
	Addr     string `mapstructure:"addr"` // Redis地址
	Password string `mapstructure:"pwd"`  // Redis密码
	Db       int    `mapstructure:"db"`   // Redis数据库
}
