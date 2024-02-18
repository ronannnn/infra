package infra

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	Mssql      = "mssql"
	Mysql      = "mysql"
	Postgresql = "postgresql"
)

type DbCfg struct {
	DbType          string `mapstructure:"db-type"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Addr            string `mapstructure:"addr"`
	Port            int    `mapstructure:"port"`
	Schema          string `mapstructure:"schema"`  // i.e. database name
	Configs         string `mapstructure:"configs"` // some extra configs appended after dsn
	MaxIdleConns    int    `mapstructure:"max-idle-conns"`
	MaxOpenConns    int    `mapstructure:"max-open-conns"`
	ConnMaxLifetime int    `mapstructure:"conn-max-lifetime"` // time unit: seconds
	EnableLog       bool   `mapstructure:"enable-log"`
}

func (dbCfg *DbCfg) dsn() (dsn string, err error) {
	switch dbCfg.DbType {
	case Mssql:
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable",
			dbCfg.Username, dbCfg.Password, dbCfg.Addr, dbCfg.Port, dbCfg.Schema)
	case Mysql:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			dbCfg.Username, dbCfg.Password, dbCfg.Addr, dbCfg.Port, dbCfg.Schema, dbCfg.Configs)
	case Postgresql:
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s",
			dbCfg.Addr, dbCfg.Port, dbCfg.Username, dbCfg.Password, dbCfg.Schema, dbCfg.Configs)
	default:
		err = fmt.Errorf("unsupported database type: %s", dbCfg.DbType)
	}
	return
}

func (dbCfg *DbCfg) dsnWithoutSchema() (dsn string, err error) {
	switch dbCfg.DbType {
	case Mssql:
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?encrypt=disable",
			dbCfg.Username, dbCfg.Password, dbCfg.Addr, dbCfg.Port)
	case Mysql:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/?%s",
			dbCfg.Username, dbCfg.Password, dbCfg.Addr, dbCfg.Port, dbCfg.Configs)
	case Postgresql:
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s %s",
			dbCfg.Addr, dbCfg.Port, dbCfg.Username, dbCfg.Password, dbCfg.Configs)
	default:
		err = fmt.Errorf("unsupported database type: %s", dbCfg.DbType)
	}
	return
}

// gormConfigs create gorm configurations
func gormConfigs(enableLog bool) *gorm.Config {
	var logMode logger.Interface
	if enableLog {
		logMode = logger.Default.LogMode(logger.Info)
	} else {
		logMode = logger.Default.LogMode(logger.Silent)
	}
	return &gorm.Config{
		Logger:                                   logMode,
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}

func NewDb(
	cfg DbCfg,
	dropTablesBeforeMigration bool,
	tables []any,
) (db *gorm.DB, err error) {
	var dbHandler typedDbHandler
	switch cfg.DbType {
	case Mssql:
		dbHandler = newMssqlHandler(cfg)
	case Mysql:
		dbHandler = newMysqlHandler(cfg)
	case Postgresql:
		dbHandler = newPostgresqlHandler(cfg)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.DbType)
	}
	if err = dbHandler.EnsureDb(); err != nil {
		return
	}
	if db, err = dbHandler.EstablishDb(); err != nil {
		return
	}
	if dropTablesBeforeMigration {
		if err = dropTables(db); err != nil {
			return
		}
	}
	err = db.AutoMigrate(tables...)
	return
}

func dropTables(db *gorm.DB) (err error) {
	var curTables []string
	if curTables, err = db.Migrator().GetTables(); err != nil {
		return
	}
	var adaptedCurTables []any
	for _, str := range curTables {
		adaptedCurTables = append(adaptedCurTables, str)
	}
	return db.Migrator().DropTable(adaptedCurTables)
}
