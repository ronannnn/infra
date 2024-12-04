package db

import (
	"fmt"
)

const (
	DbTypeMssql      = "mssql"
	DbTypeMysql      = "mysql"
	DbTypePostgresql = "postgresql"
)

type Cfg struct {
	DbType            string `mapstructure:"db-type"`
	Username          string `mapstructure:"username"`
	Password          string `mapstructure:"password"`
	Addr              string `mapstructure:"addr"`
	Port              int    `mapstructure:"port"`
	Schema            string `mapstructure:"schema"`  // i.e. database name
	Configs           string `mapstructure:"configs"` // some extra configs appended after dsn
	MaxIdleConns      int    `mapstructure:"max-idle-conns"`
	MaxOpenConns      int    `mapstructure:"max-open-conns"`
	ConnMaxLifetime   int    `mapstructure:"conn-max-lifetime"` // time unit: seconds
	EnableLog         bool   `mapstructure:"enable-log"`
	EnableAutoMigrate bool   `mapstructure:"enable-auto-migrate"`
}

func (dbCfg *Cfg) Dsn() (dsn string, err error) {
	switch dbCfg.DbType {
	case DbTypeMssql:
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable",
			dbCfg.Username, dbCfg.Password, dbCfg.Addr, dbCfg.Port, dbCfg.Schema)
	case DbTypeMysql:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			dbCfg.Username, dbCfg.Password, dbCfg.Addr, dbCfg.Port, dbCfg.Schema, dbCfg.Configs)
	case DbTypePostgresql:
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s",
			dbCfg.Addr, dbCfg.Port, dbCfg.Username, dbCfg.Password, dbCfg.Schema, dbCfg.Configs)
	default:
		err = fmt.Errorf("unsupported database type: %s", dbCfg.DbType)
	}
	return
}

func (dbCfg *Cfg) DsnWithoutSchema() (dsn string, err error) {
	switch dbCfg.DbType {
	case DbTypeMssql:
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?encrypt=disable",
			dbCfg.Username, dbCfg.Password, dbCfg.Addr, dbCfg.Port)
	case DbTypeMysql:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/?%s",
			dbCfg.Username, dbCfg.Password, dbCfg.Addr, dbCfg.Port, dbCfg.Configs)
	case DbTypePostgresql:
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s %s",
			dbCfg.Addr, dbCfg.Port, dbCfg.Username, dbCfg.Password, dbCfg.Configs)
	default:
		err = fmt.Errorf("unsupported database type: %s", dbCfg.DbType)
	}
	return
}
