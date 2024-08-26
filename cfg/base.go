package cfg

import (
	"github.com/ronannnn/infra/db"
	"github.com/ronannnn/infra/log"
	"github.com/ronannnn/infra/oss"
)

type Cfg struct {
	Sys  Sys     `mapstructure:"sys"`
	Log  log.Cfg `mapstructure:"log"`
	Auth Auth    `mapstructure:"auth"`
	Db   db.Cfg  `mapstructure:"db"`
	Dfs  oss.Cfg `mapstructure:"dfs"`
}
