package db

import (
	"fmt"

	"github.com/ronannnn/infra/cfg"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(
	dbCfg *cfg.Db,
	dropTablesBeforeMigration bool,
	tables []any,
) (db *gorm.DB, err error) {
	var dbHandler typedDbHandler
	switch dbCfg.DbType {
	case cfg.DbTypeMssql:
		dbHandler = newMssqlHandler(dbCfg)
	case cfg.DbTypeMysql:
		dbHandler = newMysqlHandler(dbCfg)
	case cfg.DbTypePostgresql:
		dbHandler = newPostgresqlHandler(dbCfg)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbCfg.DbType)
	}
	if err = dbHandler.ensureDb(); err != nil {
		return
	}
	if db, err = dbHandler.establishDb(); err != nil {
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
