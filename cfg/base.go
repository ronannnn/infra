package cfg

type Cfg struct {
	Sys Sys `mapstructure:"sys"`
	Log Log `mapstructure:"log"`
	Jwt Jwt `mapstructure:"jwt"`
	Db  Db  `mapstructure:"db"`
	Dfs Dfs `mapstructure:"dfs"`
}
