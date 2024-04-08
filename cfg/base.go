package cfg

type Cfg struct {
	Sys  Sys  `mapstructure:"sys"`
	Log  Log  `mapstructure:"log"`
	Auth Auth `mapstructure:"auth"`
	Db   Db   `mapstructure:"db"`
	Dfs  Dfs  `mapstructure:"dfs"`
}
