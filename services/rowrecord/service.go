package department

import (
	"github.com/ronannnn/infra/models"
	"gorm.io/gorm"
)

type Repo interface {
	Create(*gorm.DB, []*models.RowRecord) error
	DeleteById(*gorm.DB, uint) error
	DeleteByIds(*gorm.DB, []uint) error
	GetByTableNameAndRowId(tx *gorm.DB, tableName string, rowId uint) ([]models.RowRecord, error)
}
