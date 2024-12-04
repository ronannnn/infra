package db

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(
	dbCfg *Cfg,
	dropTablesBeforeMigration bool,
	tables []any,
) (db *gorm.DB, err error) {
	var dbHandler typedDbHandler
	switch dbCfg.DbType {
	case DbTypeMssql:
		dbHandler = newMssqlHandler(dbCfg)
	case DbTypeMysql:
		dbHandler = newMysqlHandler(dbCfg)
	case DbTypePostgresql:
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
	if dbCfg.EnableAutoMigrate {
		err = db.AutoMigrate(tables...)
	}
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
