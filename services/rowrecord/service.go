package department

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services"
	"gorm.io/gorm"
)

type Repo interface {
	services.CrudRepo[*models.RowRecord]
	GetByTableNameAndRowId(tx *gorm.DB, tableName string, rowId uint) (list []models.RowRecord, err error)
}
