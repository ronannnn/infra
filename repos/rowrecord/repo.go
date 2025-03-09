package rowrecord

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/repos"
	srv "github.com/ronannnn/infra/services/rowrecord"
	"gorm.io/gorm"
)

func New() srv.Repo {
	return &repo{}
}

type repo struct {
	repos.DefaultCrudRepo[*models.RowRecord]
}

func (s *repo) GetByTableNameAndRowId(tx *gorm.DB, tableName string, rowId uint) (list []models.RowRecord, err error) {
	err = tx.Where(&models.RowRecord{Table: tableName, RowId: rowId}).Find(&list).Error
	return
}
