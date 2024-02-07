package db

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlHandler struct {
	cfg Cfg
}

func newMysqlHandler(cfg Cfg) *mysqlHandler {
	return &mysqlHandler{cfg: cfg}
}

func (h *mysqlHandler) EnsureDb() (err error) {
	// connect without database
	var dsn string
	if dsn, err = h.cfg.DsnWithoutSchema(); err != nil {
		return
	}
	var tmpDb *gorm.DB
	if tmpDb, err = gorm.Open(
		mysql.New(mysql.Config{
			DSN:                       dsn, // DSN data source name
			DefaultStringSize:         191,
			SkipInitializeWithVersion: true, // configure according to the version automatically
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
	_, err = sqlDb.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", h.cfg.Schema))
	return
}

func (h *mysqlHandler) EstablishDb() (db *gorm.DB, err error) {
	var dsn string
	if dsn, err = h.cfg.Dsn(); err != nil {
		return
	}
	if db, err = gorm.Open(
		mysql.New(mysql.Config{
			DSN:                       dsn, // DSN data source name
			DefaultStringSize:         191,
			SkipInitializeWithVersion: true, // configure according to the version automatically
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
