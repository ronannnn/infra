package infra

import (
	"github.com/ronannnn/veken-infra/internal/db"
	"github.com/ronannnn/veken-infra/internal/dfs"
	"github.com/ronannnn/veken-infra/internal/i18n"
	"github.com/ronannnn/veken-infra/internal/jwt"
	"github.com/ronannnn/veken-infra/internal/log"
	"github.com/ronannnn/veken-infra/internal/sys"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DbCfg = db.Cfg

func NewDb(cfg DbCfg, dropTablesBeforeMigration bool, tables []any) (*gorm.DB, error) {
	return db.New(cfg, dropTablesBeforeMigration, tables)
}

type DfsCfg = dfs.Cfg

type Dfs = dfs.Dfs

func NewDfs(cfg DfsCfg) (Dfs, error) {
	return dfs.New(cfg)
}

type I18nCfg = i18n.Cfg

type I18n = i18n.I18n

func NewI18n(cfg I18nCfg) I18n {
	return i18n.New(cfg)
}

type JwtCfg = jwt.Cfg

type LogCfg = log.Cfg

func NewLog(cfg LogCfg) (*zap.SugaredLogger, error) {
	return log.New(cfg)
}

type SysCfg = sys.Cfg
