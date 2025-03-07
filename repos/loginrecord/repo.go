package loginrecord

import (
	"github.com/ronannnn/infra/models"
	srv "github.com/ronannnn/infra/services/loginrecord"
	"gorm.io/gorm"
)

func New() srv.Repo {
	return &repo{}
}

type repo struct {
}

func (s repo) Create(tx *gorm.DB, model *models.LoginRecord) error {
	return tx.Create(model).Error
}
