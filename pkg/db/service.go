package db

import (
	"fmt"

	"gorm.io/gorm"
)

type typedDbHandler interface {
	// EnsureDb ensure that database exists, create if not exists
	EnsureDb() error
	// EstablishDb connect db and create gorm.DB instance
	EstablishDb() (*gorm.DB, error)
}

func New(
	cfg *Cfg,
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
