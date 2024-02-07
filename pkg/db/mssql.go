package db

import (
	"gorm.io/gorm"
)

type mssqlHandler struct {
	cfg *Cfg
}

func newMssqlHandler(cfg *Cfg) *mssqlHandler {
	return &mssqlHandler{cfg: cfg}
}

func (h *mssqlHandler) EnsureDb() error {
	//TODO implement me
	panic("implement me")
}

func (h *mssqlHandler) EstablishDb() (*gorm.DB, error) {
	//TODO implement me
	panic("implement me")
}
