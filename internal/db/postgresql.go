package db

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresqlHandler struct {
	cfg Cfg
}

func newPostgresqlHandler(cfg Cfg) *postgresqlHandler {
	return &postgresqlHandler{cfg: cfg}
}

func (h *postgresqlHandler) EnsureDb() (err error) {
	// connect without database
	var dsn string
	if dsn, err = h.cfg.DsnWithoutSchema(); err != nil {
		return
	}
	var tmpDb *gorm.DB
	if tmpDb, err = gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  dsn, // DSN data source name
			PreferSimpleProtocol: false,
		}),
		gormConfigs(h.cfg.EnableLog),
	); err != nil {
		return
	}
	var sqlDb *sql.DB
	if sqlDb, err = tmpDb.DB(); err != nil {
		return
	}
	defer func(sqlDb *sql.DB) {
		_ = sqlDb.Close()
	}(sqlDb)
	_, err = sqlDb.Exec(fmt.Sprintf("CREATE DATABASE %s;", h.cfg.Schema))
	return
}

func (h *postgresqlHandler) EstablishDb() (db *gorm.DB, err error) {
	var dsn string
	if dsn, err = h.cfg.Dsn(); err != nil {
		return
	}
	if db, err = gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  dsn, // DSN data source name
			PreferSimpleProtocol: false,
		}),
		gormConfigs(h.cfg.EnableLog),
	); err != nil {
		return
	}
	// set further Conf
	var sqlDb *sql.DB
	if sqlDb, err = db.DB(); err != nil {
		return
	}
	// reference: https://colobu.com/2020/05/18/configuring-sql-DB-for-better-performance-2020/
	sqlDb.SetMaxOpenConns(h.cfg.MaxOpenConns)
	sqlDb.SetMaxIdleConns(h.cfg.MaxIdleConns)
	sqlDb.SetConnMaxLifetime(time.Second * time.Duration(h.cfg.ConnMaxLifetime))
	return
}
